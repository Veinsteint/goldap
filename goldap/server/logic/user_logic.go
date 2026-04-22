package logic

import (
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type UserLogic struct{}

// Add creates a new user
func (l UserLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("username already exists"))
	}
	
	if r.Mail != "" && isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorError(fmt.Errorf("email already exists"))
	}

	// Decrypt RSA encrypted password
	if r.Password != "" {
		decodeData, err := tools.RSADecrypt([]byte(r.Password), config.Conf.System.RSAPrivateBytes)
		if err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("password decryption failed"))
		}
		r.Password = string(decodeData)
		if len(r.Password) < 6 {
			return nil, tools.NewValidatorError(fmt.Errorf("password must be at least 6 characters"))
		}
		if r.ConfirmPassword != "" {
			confirmDecodeData, err := tools.RSADecrypt([]byte(r.ConfirmPassword), config.Conf.System.RSAPrivateBytes)
			if err != nil {
				return nil, tools.NewValidatorError(fmt.Errorf("confirm password decryption failed"))
			}
			if r.Password != string(confirmDecodeData) {
				return nil, tools.NewValidatorError(fmt.Errorf("passwords do not match"))
			}
		}
	} else {
		r.Password = config.Conf.Ldap.UserInitPassword
	}

	currentRoleSortMin, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get current user role"))
	}

	if len(r.RoleIds) == 0 {
		r.RoleIds = []uint{2} // Default user role
	}

	roles, err := isql.Role.GetRolesByIds(r.RoleIds)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get roles"))
	}

	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts).(int))

	// Role hierarchy check (admin bypasses)
	if currentRoleSortMin != 1 {
		if currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorError(fmt.Errorf("cannot create user with higher or equal role"))
		}
	}

	user := model.User{
		Username:      r.Username,
		Password:      r.Password,
		Nickname:      r.Nickname,
		GivenName:     r.GivenName,
		Mail:          r.Mail,
		JobNumber:     r.JobNumber,
		Mobile:        r.Mobile,
		Avatar:        r.Avatar,
		PostalAddress: r.PostalAddress,
		Departments:   r.Departments,
		Position:      r.Position,
		Introduction:  r.Introduction,
		Status:        r.Status,
		Creator:       ctxUser.Username,
		DepartmentId:  tools.SliceToString(r.DepartmentId, ","),
		Source:        r.Source,
		Roles:         roles,
		UIDNumber:     r.UIDNumber,
		GIDNumber:     r.GIDNumber,
		HomeDirectory: r.HomeDirectory,
		LoginShell:    r.LoginShell,
		Gecos:         r.Gecos,
	}

	if user.Source == "" {
		user.Source = "platform"
	}

	var groups []*model.Group
	if user.DepartmentId != "" {
		groups, err = isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %s", err.Error()))
		}
	}

	err = CommonAddUser(&user, groups)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "validator") {
			return nil, err
		}
		return nil, tools.NewOperationError(fmt.Errorf("failed to add user: %s", err.Error()))
	}
	return nil, nil
}

// List returns paginated user list
func (l UserLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	users, err := isql.User.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user list: %s", err.Error()))
	}

	rets := make([]model.User, 0)
	for _, user := range users {
		rets = append(rets, *user)
	}
	count, err := isql.User.ListCount(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user count: %s", err.Error()))
	}

	return response.UserListRsp{
		Total: int(count),
		Users: rets,
	}, nil
}

// Update modifies user data
func (l UserLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if !isql.User.Exist(tools.H{"id": r.ID}) {
		return nil, tools.NewMySqlError(fmt.Errorf("user not found"))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	var currentRoleSorts []int
	for _, role := range ctxUser.Roles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}

	var reqRoleSorts []int
	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get roles"))
	}
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}

	currentRoleSortMin := funk.MinInt(currentRoleSorts).(int)
	reqRoleSortMin := funk.MinInt(reqRoleSorts).(int)

	// Role hierarchy check (admin bypasses)
	if currentRoleSortMin != 1 {
		if int(r.ID) == int(ctxUser.ID) {
			reqDiff, currentDiff := funk.Difference(reqRoleSorts, currentRoleSorts)
			if len(reqDiff.([]int)) > 0 || len(currentDiff.([]int)) > 0 {
				return nil, tools.NewValidatorError(fmt.Errorf("cannot modify own roles"))
			}
		}

		minRoleSorts, err := isql.User.GetUserMinRoleSortsByIds([]uint{uint(r.ID)})
		if err != nil || len(minRoleSorts) == 0 {
			return nil, tools.NewValidatorError(fmt.Errorf("failed to get user role"))
		}
		if currentRoleSortMin >= minRoleSorts[0] || currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorError(fmt.Errorf("cannot update user with higher or equal role"))
		}
	}

	oldData := new(model.User)
	err = isql.User.Find(tools.H{"id": r.ID}, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	// Filter invalid department selections
	var depts string
	var deptids []uint
	for _, v := range strings.Split(r.Departments, ",") {
		if v != "请选择部门信息" {
			depts += v + ","
		}
	}
	for _, j := range r.DepartmentId {
		if j != 0 {
			deptids = append(deptids, j)
		}
	}

	user := *oldData
	user.Username = r.Username
	user.Nickname = r.Nickname
	user.GivenName = r.GivenName
	user.Mail = r.Mail
	user.JobNumber = r.JobNumber
	user.Mobile = r.Mobile
	user.Avatar = r.Avatar
	user.PostalAddress = r.PostalAddress
	user.Departments = depts
	user.Position = r.Position
	user.Introduction = r.Introduction
	user.Creator = ctxUser.Username
	user.DepartmentId = tools.SliceToString(deptids, ",")
	user.Roles = roles

	// Handle Unix attributes
	if r.UIDNumber > 0 {
		user.UIDNumber = r.UIDNumber
	} else if oldData.UIDNumber == 0 {
		uid, err := isql.User.GetNextUIDNumber()
		if err == nil {
			user.UIDNumber = uid
		}
	}

	if r.GIDNumber > 0 {
		user.GIDNumber = r.GIDNumber
	} else if oldData.GIDNumber == 0 && user.UIDNumber > 0 {
		user.GIDNumber = user.UIDNumber
	}

	if user.UIDNumber > 0 && user.GIDNumber == 0 {
		user.GIDNumber = user.UIDNumber
	}

	if r.HomeDirectory != "" {
		user.HomeDirectory = r.HomeDirectory
	} else if oldData.HomeDirectory == "" {
		user.HomeDirectory = fmt.Sprintf("/home/%s", user.Username)
	}

	if r.LoginShell != "" {
		user.LoginShell = r.LoginShell
	} else if oldData.LoginShell == "" {
		user.LoginShell = "/bin/bash"
	}

	if r.Gecos != "" {
		user.Gecos = r.Gecos
	} else if oldData.Gecos == "" {
		if user.Nickname != "" {
			user.Gecos = user.Nickname
		} else {
			user.Gecos = user.Username
		}
	}

	user.Source = oldData.Source
	user.UserDN = oldData.UserDN
	user.SyncState = oldData.SyncState

	if err = CommonUpdateUser(oldData, &user, r.DepartmentId); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to update user: %s", err.Error()))
	}

	// Sync user data to pre-config if exists
	_ = isql.UserPreConfig.SyncFromUser(&user)

	return nil, nil
}

// Delete removes users
func (l UserLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.UserIds {
		filter := tools.H{"id": int(id)}
		if !isql.User.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("user not found"))
		}
	}

	roleMinSortList, err := isql.User.GetUserMinRoleSortsByIds(r.UserIds)
	if err != nil || len(roleMinSortList) == 0 {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get user roles"))
	}

	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get current user role"))
	}

	if funk.Contains(r.UserIds, ctxUser.ID) {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot delete self"))
	}

	for _, sort := range roleMinSortList {
		if int(minSort) > sort {
			return nil, tools.NewValidatorError(fmt.Errorf("cannot delete user with higher role"))
		}
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %s", err.Error()))
	}

	// Delete from LDAP first
	for _, user := range users {
		if user.UserDN != "" {
			err := ildap.User.Delete(user.UserDN)
			if err != nil {
				common.Log.Errorf("Failed to delete LDAP user [%s]: %v", user.Username, err)
			}
		}
	}

	// Delete from MySQL
	err = isql.User.Delete(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete MySQL users: %s", err.Error()))
	}

	return nil, nil
}

// ChangePwd changes user password
func (l UserLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	decodeOldPassword, err := tools.RSADecrypt([]byte(r.OldPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("old password decryption failed"))
	}
	decodeNewPassword, err := tools.RSADecrypt([]byte(r.NewPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("new password decryption failed"))
	}
	r.OldPassword = string(decodeOldPassword)
	r.NewPassword = string(decodeNewPassword)

	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	if tools.NewParPasswd(user.Password) != r.OldPassword {
		return nil, tools.NewValidatorError(fmt.Errorf("incorrect old password"))
	}

	err = ildap.User.ChangePwd(user.UserDN, "", r.NewPassword)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to update LDAP password: %s", err.Error()))
	}

	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(r.NewPassword))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update MySQL password: %s", err.Error()))
	}

	return nil, nil
}

// ResetPassword resets user password and sends email notification
func (l UserLogic) ResetPassword(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserResetPasswordReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if !isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not found"))
	}

	user := new(model.User)
	err := isql.User.Find(tools.H{"username": r.Username}, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user: %s", err.Error()))
	}

	newPassword := tools.GenerateRandomPassword()

	err = ildap.User.ChangePwd(user.UserDN, "", newPassword)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to update LDAP password: %s", err.Error()))
	}

	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newPassword))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update MySQL password: %s", err.Error()))
	}

	if err := tools.SendPasswordResetNotification(user.Username, user.Nickname, user.Mail, newPassword); err != nil {
		common.Log.Warnf("Failed to send password reset email: %s, %v", user.Username, err)
	}

	return map[string]string{"newPassword": newPassword}, nil
}

// ChangeUserStatus changes user active/inactive status
func (l UserLogic) ChangeUserStatus(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangeUserStatusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.ID}
	if !isql.User.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not found"))
	}

	user := new(model.User)
	err := isql.User.Find(filter, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user: %s", err.Error()))
	}

	if r.Status == user.Status {
		if r.Status == 1 {
			return nil, tools.NewValidatorError(fmt.Errorf("user already active"))
		}
		return nil, tools.NewValidatorError(fmt.Errorf("user already inactive"))
	}

	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("failed to get current user role"))
	}

	if int(minSort) != 1 {
		return nil, tools.NewValidatorError(fmt.Errorf("only admin can change user status"))
	}

	if r.Status == 2 {
		err = ildap.User.Delete(user.UserDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("failed to delete LDAP user: %s", err.Error()))
		}
	} else {
		err = ildap.User.Add(user)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("failed to add LDAP user: %s", err.Error()))
		}
	}

	err = isql.User.ChangeStatus(int(r.ID), int(r.Status))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update user status: %s", err.Error()))
	}
	return nil, nil
}

// GetUserInfo returns current logged-in user info
func (l UserLogic) GetUserInfo(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserGetUserInfoReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	_ = r

	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user: %s", err.Error()))
	}
	return user, nil
}
