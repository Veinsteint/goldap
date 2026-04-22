package logic

import (
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"
	
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

var (
	ReqAssertErr = tools.NewRspError(tools.SystemErr, fmt.Errorf("request assertion failed"))

	Api                   = &ApiLogic{}
	User                  = &UserLogic{}
	Group                 = &GroupLogic{}
	Role                  = &RoleLogic{}
	Menu                  = &MenuLogic{}
	OperationLog          = &OperationLogLogic{}
	OpenLdap              = &OpenLdapLogic{}
	Sql                   = &SqlLogic{}
	Base                  = &BaseLogic{}
	FieldRelation         = &FieldRelationLogic{}
	SSHKey                = &SSHKeyLogic{}
	IPGroup               = &IPGroupLogic{}
	IPGroupUserPermission = &IPGroupUserPermissionLogic{}
	GroupUserPermission   = &GroupUserPermissionLogic{}
	SudoRule              = &SudoRuleLogic{}
	PendingUser           = &PendingUserLogic{}
	UserPreConfig         = &UserPreConfigLogic{}
	SystemConfig          = &SystemConfigLogic{}

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// CommonAddGroup creates group in LDAP and MySQL
func CommonAddGroup(group *model.Group) error {
	err := ildap.Group.Add(group)
	if err != nil {
		return err
	}

	err = isql.Group.Add(group)
	if err != nil {
		return err
	}

	return nil
}

// CommonUpdateGroup updates group in both LDAP and MySQL
func CommonUpdateGroup(oldGroup, newGroup *model.Group) error {
	if !config.Conf.Ldap.GroupNameModify {
		newGroup.GroupName = oldGroup.GroupName
	}

	err := ildap.Group.Update(oldGroup, newGroup)
	if err != nil {
		return err
	}
	err = isql.Group.Update(newGroup)
	if err != nil {
		return err
	}
	return nil
}

// CommonAddUser creates user with group assignment
func CommonAddUser(user *model.User, groups []*model.Group) error {
	// Set default values
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	if user.GivenName == "" {
		user.GivenName = user.Nickname
	}
	if user.Introduction == "" {
		user.Introduction = user.Nickname
	}
	if user.Mail == "" {
		if len(config.Conf.Ldap.DefaultEmailSuffix) > 0 {
			user.Mail = user.Username + "@" + config.Conf.Ldap.DefaultEmailSuffix
		} else {
			user.Mail = user.Username + "@163.com"
		}
	}
	if user.Departments == "" {
		user.Departments = "CMPLAB"
	}
	if user.Position == "" {
		user.Position = "CIMR"
	}

	// Determine target group (prefer OU type)
	var targetGroup *model.Group
	if len(groups) > 0 {
		for _, g := range groups {
			if strings.HasPrefix(g.GroupDN, "ou=") {
				targetGroup = g
				break
			}
		}
		if targetGroup == nil {
			targetGroup = groups[0]
		}
	} else {
		// Use default group CMPLabHPC
		defaultGroup := new(model.Group)
		err := isql.Group.Find(tools.H{"group_name": "CMPLabHPC"}, defaultGroup)
		if err != nil {
			defaultGroup = &model.Group{
				GroupName: "CMPLabHPC",
				Remark:    "Default group",
				GroupType: "ou",
				GroupDN:   fmt.Sprintf("ou=CMPLabHPC,%s", config.Conf.Ldap.BaseDN),
				Creator:   "system",
				Source:    "platform",
			}
			if err := CommonAddGroup(defaultGroup); err != nil {
				common.Log.Warnf("Failed to create default group CMPLabHPC: %v", err)
				user.UserDN = fmt.Sprintf("uid=%s,ou=CMPLabHPC,%s", user.Username, config.Conf.Ldap.BaseDN)
			} else {
				targetGroup = defaultGroup
			}
		} else {
			targetGroup = defaultGroup
		}
	}

	// Set UserDN based on group
	if targetGroup != nil && strings.HasPrefix(targetGroup.GroupDN, "ou=") {
		user.UserDN = fmt.Sprintf("uid=%s,%s", user.Username, targetGroup.GroupDN)
	} else {
		user.UserDN = fmt.Sprintf("uid=%s,ou=CMPLabHPC,%s", user.Username, config.Conf.Ldap.BaseDN)
	}

	// Add to MySQL (GORM hooks sync to LDAP)
	err := isql.User.Add(user)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "Duplicate entry") {
			if strings.Contains(errStr, "users.username") {
				return tools.NewValidatorError(fmt.Errorf("username already exists"))
			}
		}
		return tools.NewMySqlError(fmt.Errorf("failed to create user in MySQL: %s", err.Error()))
	}

	// Send notification email
	if err := tools.SendUserCreationNotification(user.Username, user.Nickname, user.Mail, tools.NewParPasswd(user.Password)); err != nil {
		common.Log.Warnf("Failed to send user creation email: %s, %v", user.Username, err)
	}

	// Add user to groups
	for _, group := range groups {
		err := isql.Group.AddUserToGroup(group, []model.User{*user})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to add user to group: %s", err.Error()))
		}
		
		// Sync to LDAP for non-OU groups (posixGroup like docker, groupOfUniqueNames, etc.)
		isOUGroup := len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou="
		if !isOUGroup {
			err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
			if err != nil {
				// Rollback
				_ = isql.Group.RemoveUserFromGroup(group, []model.User{*user})
				common.Log.Warnf("Failed to add user to LDAP group [%s]: %v", group.GroupDN, err)
			}
		}
	}

	return nil
}

// CommonUpdateUser updates user in LDAP and MySQL with group changes
func CommonUpdateUser(oldUser, newUser *model.User, groupId []uint) error {
	if !config.Conf.Ldap.UserNameModify {
		newUser.Username = oldUser.Username
	}

	err := ildap.User.Update(oldUser.Username, newUser)
	if err != nil {
		return tools.NewLdapError(fmt.Errorf("failed to update LDAP user: %s", err.Error()))
	}

	err = isql.User.Update(newUser, true)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("failed to update MySQL user: %s", err.Error()))
	}

	// Handle group changes
	oldDeptIds := tools.StringToSlice(oldUser.DepartmentId, ",")
	addDeptIds, removeDeptIds := tools.ArrUintCmp(oldDeptIds, groupId)

	// Add to new groups
	addgroups, err := isql.Group.GetGroupByIds(addDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
	}
	for _, group := range addgroups {
		if len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou=" {
			continue
		}
		err := isql.Group.AddUserToGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to add user to group: %s", err.Error()))
		}
		err = ildap.Group.AddUserToGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewLdapError(fmt.Errorf("failed to add user to LDAP group: %s", err.Error()))
		}
	}

	// Remove from old groups
	removegroups, err := isql.Group.GetGroupByIds(removeDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
	}
	for _, group := range removegroups {
		if len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou=" {
			continue
		}
		err := isql.Group.RemoveUserFromGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to remove user from group: %s", err.Error()))
		}
		err = ildap.Group.RemoveUserFromGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("failed to remove user from LDAP group: %s", err.Error()))
		}
	}
	return nil
}

// BuildGroupData builds group model from remote data using field mapping
func BuildGroupData(flag string, remoteData map[string]any) (*model.Group, error) {
	output, err := json.Marshal(&remoteData)
	if err != nil {
		return nil, err
	}

	oldData := new(model.FieldRelation)
	err = isql.FieldRelation.Find(tools.H{"flag": flag + "_group"}, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}
	frs, err := tools.JsonToMap(string(oldData.Attributes))
	if err != nil {
		return nil, tools.NewOperationError(err)
	}

	g := &model.Group{}
	for system, remote := range frs {
		switch system {
		case "groupName":
			g.SetGroupName(gjson.Get(string(output), remote).String())
		case "remark":
			g.SetRemark(gjson.Get(string(output), remote).String())
		case "sourceDeptId":
			g.SetSourceDeptId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		case "sourceDeptParentId":
			g.SetSourceDeptParentId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		}
	}
	return g, nil
}

// BuildUserData builds user model from remote data using field mapping
func BuildUserData(flag string, remoteData map[string]any) (*model.User, error) {
	output, err := json.Marshal(&remoteData)
	if err != nil {
		return nil, err
	}

	fieldRelationSource := new(model.FieldRelation)
	err = isql.FieldRelation.Find(tools.H{"flag": flag + "_user"}, fieldRelationSource)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}
	fieldRelation, err := tools.JsonToMap(string(fieldRelationSource.Attributes))
	if err != nil {
		return nil, tools.NewOperationError(err)
	}

	name := gjson.Get(string(output), fieldRelation["username"]).String()
	if len(name) == 0 {
		common.Log.Warnf("User missing username: %s", output)
		return nil, nil
	}

	u := &model.User{}
	for system, remote := range fieldRelation {
		switch system {
		case "username":
			u.SetUserName(gjson.Get(string(output), remote).String())
		case "nickname":
			u.SetNickName(gjson.Get(string(output), remote).String())
		case "givenName":
			u.SetGivenName(gjson.Get(string(output), remote).String())
		case "mail":
			u.SetMail(gjson.Get(string(output), remote).String())
		case "jobNumber":
			u.SetJobNumber(gjson.Get(string(output), remote).String())
		case "mobile":
			u.SetMobile(gjson.Get(string(output), remote).String())
		case "avatar":
			u.SetAvatar(gjson.Get(string(output), remote).String())
		case "postalAddress":
			u.SetPostalAddress(gjson.Get(string(output), remote).String())
		case "position":
			u.SetPosition(gjson.Get(string(output), remote).String())
		case "introduction":
			u.SetIntroduction(gjson.Get(string(output), remote).String())
		case "sourceUserId":
			u.SetSourceUserId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		case "sourceUnionId":
			u.SetSourceUnionId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		}
	}
	return u, nil
}

// ConvertDeptData converts department data to group models
func ConvertDeptData(flag string, remoteData []map[string]any) (groups []*model.Group, err error) {
	for _, dept := range remoteData {
		group, err := BuildGroupData(flag, dept)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return
}

// ConvertUserData converts user data to user models
func ConvertUserData(flag string, remoteData []map[string]any) (users []*model.User, err error) {
	for _, staff := range remoteData {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff["department_ids"].([]string))
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to convert dept IDs: %s", err.Error()))
		}
		user, err := BuildUserData(flag, staff)
		if err != nil {
			return nil, err
		}
		if user != nil {
			user.DepartmentId = tools.SliceToString(groupIds, ",")
			users = append(users, user)
		}
	}
	return
}

// InitCron starts scheduled sync tasks
func InitCron() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 */2 * * * *", func() {
		_ = SearchGroupDiff()
		_ = SearchUserDiff()
	})
	if err != nil {
		common.Log.Errorf("Failed to start sync cron job: %v", err)
	}
	c.Start()
}

// GroupListToTree converts flat group list to tree structure
func GroupListToTree(rootId string, groupList []*model.Group) *model.Group {
	rootGroup := &model.Group{SourceDeptId: rootId}
	rootGroup.Children = groupListToTree(rootGroup, groupList)
	return rootGroup
}

func groupListToTree(rootGroup *model.Group, list []*model.Group) []*model.Group {
	children := make([]*model.Group, 0)
	for _, group := range list {
		if group.SourceDeptParentId == rootGroup.SourceDeptId {
			children = append(children, group)
		}
	}
	for _, group := range children {
		group.Children = groupListToTree(group, list)
	}
	return children
}
