package logic

import (
	"fmt"
	"strconv"
	"strings"

	"goldap-server/config"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/gin-gonic/gin"
)

type GroupLogic struct{}

// Add creates a new group
func (l GroupLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	group := model.Group{
		GroupType: r.GroupType,
		ParentId:  r.ParentId,
		GroupName: r.GroupName,
		GIDNumber: r.GIDNumber,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
		Source:    "platform",
	}

	dnPrefix := r.GroupType
	if r.GroupType == "posix" {
		dnPrefix = "cn"
	}

	if r.ParentId == 0 {
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = "platform_0"
		group.GroupDN = fmt.Sprintf("%s=%s,%s", dnPrefix, r.GroupName, config.Conf.Ldap.BaseDN)
	} else {
		parentGroup := new(model.Group)
		err := isql.Group.Find(tools.H{"id": r.ParentId}, parentGroup)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("parent group not found"))
		}
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = fmt.Sprintf("%s_%d", parentGroup.Source, r.ParentId)
		group.GroupDN = fmt.Sprintf("%s=%s,%s", dnPrefix, r.GroupName, parentGroup.GroupDN)
	}

	if isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		return nil, tools.NewValidatorError(fmt.Errorf("group DN already exists"))
	}

	err = isql.Group.Add(&group)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to create MySQL group"))
	}
	
	if group.GroupType == "cn" && !strings.Contains(group.GroupDN, "ou=sudoers,") {
		if len(group.Users) == 0 {
			_ = isql.Group.ChangeSyncState(int(group.ID), 2)
			common.Log.Infof("Group %s created in MySQL. LDAP group will be created when first user is added.", group.GroupName)
		} else {
			// Ensure LDAP group exists
			err = ildap.Group.Add(&group)
			if err != nil {
				if strings.Contains(err.Error(), "uniqueMember") || strings.Contains(err.Error(), "Object Class Violation") {
					_ = isql.Group.ChangeSyncState(int(group.ID), 2)
					common.Log.Infof("Group %s created in MySQL. LDAP group will be created when first user is added.", group.GroupName)
				} else if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
					// Group already exists
					_ = isql.Group.ChangeSyncState(int(group.ID), 1)
				} else {
					_ = isql.Group.ChangeSyncState(int(group.ID), 2)
					common.Log.Warnf("Failed to sync group %s to LDAP: %v", group.GroupName, err)
				}
			} else {
				_ = isql.Group.ChangeSyncState(int(group.ID), 1)
			}
		}
	} else {
		if err := ildap.Group.Add(&group); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				_ = isql.Group.ChangeSyncState(int(group.ID), 1)
			}
		}
	}

	return nil, nil
}

// List returns paginated group list
func (l GroupLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	groups, err := isql.Group.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
	}

	rets := make([]model.Group, 0)
	for _, group := range groups {
		rets = append(rets, *group)
	}
	count, err := isql.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get group count"))
	}

	return response.GroupListRsp{
		Total:  count,
		Groups: rets,
	}, nil
}

// GetTree returns group tree structure
func (l GroupLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	groups, err := isql.Group.ListTree(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
	}

	tree := isql.GenGroupTree(0, groups)
	return tree, nil
}

// Update modifies group data
func (l GroupLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	oldGroup := new(model.Group)
	err = isql.Group.Find(filter, oldGroup)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	newGroup := *oldGroup
	newGroup.Remark = r.Remark
	newGroup.Creator = ctxUser.Username
	newGroup.IPRanges = r.IPRanges
	// Update GIDNumber for posixGroup
	if oldGroup.GroupType == "posix" && r.GIDNumber > 0 {
		newGroup.GIDNumber = r.GIDNumber
	}

	if !config.Conf.Ldap.GroupNameModify {
		newGroup.GroupName = oldGroup.GroupName
	} else if r.GroupName != oldGroup.GroupName {
		newGroup.GroupName = r.GroupName
		if oldGroup.ParentId == 0 {
			newGroup.GroupDN = fmt.Sprintf("%s=%s,%s", oldGroup.GroupType, r.GroupName, config.Conf.Ldap.BaseDN)
		} else {
			parentGroup := new(model.Group)
			err := isql.Group.Find(tools.H{"id": oldGroup.ParentId}, parentGroup)
			if err != nil {
				return nil, tools.NewMySqlError(fmt.Errorf("parent group not found"))
			}
			newGroup.GroupDN = fmt.Sprintf("%s=%s,%s", oldGroup.GroupType, r.GroupName, parentGroup.GroupDN)
		}
	} else {
		newGroup.GroupName = r.GroupName
	}

	err = ildap.Group.Update(oldGroup, &newGroup)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to update LDAP group: %s", err.Error()))
	}
	err = isql.Group.Update(&newGroup)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to update MySQL group"))
	}
	return nil, nil
}

// Delete removes groups
func (l GroupLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !isql.Group.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
		}
	}

	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
	}

	for _, group := range groups {
		// Prevent deletion of default sudoers groups
		if group.GroupName == "sudouser-nopasswd" || group.GroupName == "sudouser-other" || group.GroupName == "sudoers" {
			return nil, tools.NewMySqlError(fmt.Errorf("cannot delete default sudoers group: %s", group.GroupName))
		}
		if strings.Contains(group.GroupDN, "cn=sudouser-nopasswd,ou=sudoers,") || 
		   strings.Contains(group.GroupDN, "cn=sudouser-other,ou=sudoers,") ||
		   (strings.Contains(group.GroupDN, "ou=sudoers,") && group.GroupType == "ou") {
			return nil, tools.NewMySqlError(fmt.Errorf("cannot delete default sudoers group: %s", group.GroupName))
		}

		// Check for child groups
		filter := tools.H{"parent_id": int(group.ID)}
		if isql.Group.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("delete child groups first"))
		}

		err = ildap.Group.Delete(group.GroupDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("failed to delete LDAP group: %s", err.Error()))
		}
	}

	err = isql.Group.Delete(groups)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete groups: %s", err.Error()))
	}

	return nil, nil
}

// AddUser adds users to group
func (l GroupLogic) AddUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %s", err.Error()))
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get group: %s", err.Error()))
	}

	// OU groups don't directly contain users as members (organized via DN hierarchy)
	isOUGroup := len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou="
	
	err = isql.Group.AddUserToGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add users to group: %s", err.Error()))
	}

	// Add to LDAP for non-OU groups
	if !isOUGroup {
		for _, user := range users {
			err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
			if err != nil {
				_ = isql.Group.RemoveUserFromGroup(group, users)
				return nil, tools.NewLdapError(fmt.Errorf("failed to add user to LDAP group: %s", err.Error()))
			}
		}
	}

	for _, user := range users {
		oldData := new(model.User)
		err = isql.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		newData := oldData
		// Update DepartmentId - avoid duplicate and handle empty strings
		var deptIds []string
		if oldData.DepartmentId != "" {
			deptIds = strings.Split(oldData.DepartmentId, ",")
		}
		groupIDStr := strconv.Itoa(int(r.GroupID))
		exists := false
		for _, id := range deptIds {
			if id == groupIDStr {
				exists = true
				break
			}
		}
		if !exists {
			deptIds = append(deptIds, groupIDStr)
		}
		newData.DepartmentId = strings.Join(deptIds, ",")
		
		// Update Departments - avoid duplicate and handle empty strings
		var depts []string
		if oldData.Departments != "" {
			depts = strings.Split(oldData.Departments, ",")
		}
		exists = false
		for _, dept := range depts {
			if dept == group.GroupName {
				exists = true
				break
			}
		}
		if !exists {
			depts = append(depts, group.GroupName)
		}
		newData.Departments = strings.Join(depts, ",")
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("failed to update user departments: %s", err.Error()))
		}
	}

	return nil, nil
}

func (l GroupLogic) updataUser(newUser *model.User) error {
	err := isql.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("failed to update MySQL user: %s", err.Error()))
	}
	return nil
}

// RemoveUser removes users from group
func (l GroupLogic) RemoveUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupRemoveUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %s", err.Error()))
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get group: %s", err.Error()))
	}

	isOUGroup := len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou="
	
	// Remove from LDAP for non-OU groups
	if !isOUGroup {
		for _, user := range users {
			err := ildap.Group.RemoveUserFromGroup(group.GroupDN, user.UserDN)
			if err != nil {
				return nil, tools.NewLdapError(fmt.Errorf("failed to remove user from LDAP group: %s", err.Error()))
			}
		}
	}

	err = isql.Group.RemoveUserFromGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to remove users from MySQL group: %s", err.Error()))
	}

	for _, user := range users {
		oldData := new(model.User)
		err = isql.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		
		// Create a copy to avoid modifying the original
		newData := *oldData

		var newDepts []string
		var newDeptIds []string
		groupIDStr := strconv.Itoa(int(r.GroupID))
		
		// Remove group from Departments, filter out empty strings
		if oldData.Departments != "" {
			for _, v := range strings.Split(oldData.Departments, ",") {
				v = strings.TrimSpace(v)
				if v != "" && v != group.GroupName {
					newDepts = append(newDepts, v)
				}
			}
		}
		
		// Remove group ID from DepartmentId, filter out empty strings
		if oldData.DepartmentId != "" {
			for _, v := range strings.Split(oldData.DepartmentId, ",") {
				v = strings.TrimSpace(v)
				if v != "" && v != groupIDStr {
					newDeptIds = append(newDeptIds, v)
				}
			}
		}

		newData.Departments = strings.Join(newDepts, ",")
		newData.DepartmentId = strings.Join(newDeptIds, ",")
		
		// Update user - this will sync to LDAP and update businessCategory
		err = isql.User.Update(&newData)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("failed to update user departments: %s", err.Error()))
		}
	}

	return nil, nil
}

// UserInGroup returns users in specified group
func (l GroupLogic) UserInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get group: %s", err.Error()))
	}

	rets := make([]response.Guser, 0)
	for _, user := range group.Users {
		if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
			continue
		}
		rets = append(rets, response.Guser{
			UserId:       int64(user.ID),
			UserName:     user.Username,
			NickName:     user.Nickname,
			Mail:         user.Mail,
			JobNumber:    user.JobNumber,
			Mobile:       user.Mobile,
			Introduction: user.Introduction,
		})
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}

// UserNoInGroup returns users not in specified group
func (l GroupLogic) UserNoInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserNoInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("group not found"))
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get group: %s", err.Error()))
	}

	userList, err := isql.User.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %s", err.Error()))
	}

	rets := make([]response.Guser, 0)
	for _, user := range userList {
		if user.Username == "admin" {
			continue
		}
		
		in := true
		for _, groupUser := range group.Users {
			if user.Username == groupUser.Username {
				in = false
				break
			}
		}
		if in {
			if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
				continue
			}
			rets = append(rets, response.Guser{
				UserId:       int64(user.ID),
				UserName:     user.Username,
				NickName:     user.Nickname,
				Mail:         user.Mail,
				JobNumber:    user.JobNumber,
				Mobile:       user.Mobile,
				Introduction: user.Introduction,
			})
		}
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}
