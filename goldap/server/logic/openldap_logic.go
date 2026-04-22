package logic

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/model"
	"goldap-server/public/client/openldap"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type OpenLdapLogic struct{}

// SyncOpenLdapDepts syncs departments from OpenLDAP
func (d *OpenLdapLogic) SyncOpenLdapDepts(c *gin.Context, req any) (data any, rspError any) {
	depts, err := openldap.GetAllDepts()
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to get OpenLDAP departments: %v", err))
	}
	if len(depts) == 0 {
		return nil, tools.NewOperationError(errors.New("no departments found"))
	}

	groups := make([]*model.Group, 0)
	for _, dept := range depts {
		groups = append(groups, &model.Group{
			GroupName:          dept.Name,
			Remark:             dept.Remark,
			SourceDeptId:       dept.Id,
			SourceDeptParentId: dept.ParentId,
			GroupDN:            dept.DN,
		})
	}

	deptTree := GroupListToTree("0", groups)

	if err = d.addDepts(deptTree.Children); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d OpenLdapLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		if err := d.AddDepts(dept); err != nil {
			return tools.NewOperationError(fmt.Errorf("failed to add department %s: %v", dept.GroupName, err))
		}
		if len(dept.Children) != 0 {
			if err := d.addDepts(dept.Children); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddDepts adds department data
func (d OpenLdapLogic) AddDepts(group *model.Group) error {
	if !isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		group.Creator = "system"
		group.GroupType = strings.Split(strings.Split(group.GroupDN, ",")[0], "=")[0]
		parentid, err := d.getParentGroupID(group)
		if err != nil {
			return err
		}
		group.ParentId = parentid
		group.Source = "openldap"
		return isql.Group.Add(group)
	}
	return nil
}

func (d OpenLdapLogic) getParentGroupID(group *model.Group) (id uint, err error) {
	parentGroup := new(model.Group)
	err = isql.Group.Find(tools.H{"source_dept_id": group.SourceDeptParentId}, parentGroup)
	if err != nil {
		return id, tools.NewMySqlError(fmt.Errorf("failed to find parent department: %v, %s", err, group.GroupName))
	}
	return parentGroup.ID, nil
}

// SyncOpenLdapUsers syncs users from OpenLDAP
func (d OpenLdapLogic) SyncOpenLdapUsers(c *gin.Context, req any) (data any, rspError any) {
	staffs, err := openldap.GetAllUsers()
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to get OpenLDAP users: %v", err))
	}
	if len(staffs) == 0 {
		return nil, tools.NewOperationError(errors.New("no users found"))
	}

	for _, staff := range staffs {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff.DepartmentIds)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to convert department IDs for user %s: %v", staff.Name, err))
		}

		roles, err := isql.Role.GetRolesByIds([]uint{2})
		if err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("failed to get roles for user %s: %v", staff.Name, err))
		}

		if err := d.AddUsers(&model.User{
			Username:      staff.Name,
			Nickname:      staff.DisplayName,
			GivenName:     staff.GivenName,
			Mail:          staff.Mail,
			JobNumber:     staff.EmployeeNumber,
			Mobile:        staff.Mobile,
			PostalAddress: staff.PostalAddress,
			Departments:   staff.BusinessCategory,
			Position:      staff.DepartmentNumber,
			Introduction:  staff.CN,
			Creator:       "system",
			Source:        "openldap",
			DepartmentId:  tools.SliceToString(groupIds, ","),
			SourceUserId:  staff.Name,
			SourceUnionId: staff.Name,
			Roles:         roles,
			UserDN:        staff.DN,
		}); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("failed to add user %s: %v", staff.Name, err))
		}
	}

	return nil, nil
}

// AddUsers adds user data
func (d OpenLdapLogic) AddUsers(user *model.User) error {
	if !isql.User.Exist(tools.H{"user_dn": user.UserDN}) {
		if user.Departments == "" {
			user.Departments = "Default"
		}
		if user.GivenName == "" {
			user.GivenName = user.Nickname
		}
		if user.PostalAddress == "" {
			user.PostalAddress = "Default"
		}
		if user.Position == "" {
			user.Position = "Default"
		}
		if user.Introduction == "" {
			user.Introduction = user.Nickname
		}
		if user.JobNumber == "" {
			user.JobNumber = "N/A"
		}

		if err := isql.User.Add(user); err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to create user in MySQL: %v", err))
		}

		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to get groups: %v", err))
		}
		for _, group := range groups {
			if len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou=" {
				continue
			}
			if err := isql.Group.AddUserToGroup(group, []model.User{*user}); err != nil {
				return tools.NewMySqlError(fmt.Errorf("failed to add user to group: %v", err))
			}
		}
	}
	return nil
}
