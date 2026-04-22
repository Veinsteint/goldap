package logic

import (
	"fmt"
	"time"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type PendingUserLogic struct{}

// Register submits user registration for admin approval
func (l PendingUserLogic) Register(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserRegisterReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Check if registration is disabled
	registrationDisabled := isql.SystemConfig.GetValue("registration_disabled", "false") == "true"
	if registrationDisabled {
		return nil, tools.NewValidatorError(fmt.Errorf("系统升级维护中，已禁止注册，请联系管理员"))
	}

	// Check registration mode
	regMode := "preconfig" // default mode
	if config.Conf.Registration != nil && config.Conf.Registration.Mode != "" {
		regMode = config.Conf.Registration.Mode
	}

	// if regMode is "preconfig", only allow usernames in pre-config
	if regMode == "preconfig" {
		preConfig, _ := isql.UserPreConfig.GetByUsername(r.Username)
		if preConfig == nil {
			return nil, tools.NewValidatorError(fmt.Errorf("用户名 [%s] 不在预配置列表中，无法注册。请联系管理员或查看有效用户名列表", r.Username))
		}
		if preConfig.IsUsed {
			return nil, tools.NewValidatorError(fmt.Errorf("用户名 [%s] 已被注册使用", r.Username))
		}
	}

	if isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("username already exists"))
	}
	if isql.PendingUser.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("registration already submitted, pending approval"))
	}

	if isql.User.Exist(tools.H{"mail": r.Email}) {
		return nil, tools.NewValidatorError(fmt.Errorf("email already exists"))
	}
	if isql.PendingUser.Exist(tools.H{"mail": r.Email}) {
		return nil, tools.NewValidatorError(fmt.Errorf("email already submitted, pending approval"))
	}

	decodeData, err := tools.RSADecrypt([]byte(r.Password), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("password decryption failed"))
	}
	password := string(decodeData)
	if len(password) < 6 {
		return nil, tools.NewValidatorError(fmt.Errorf("password must be at least 6 characters"))
	}

	pendingUser := &model.PendingUser{
		Username: r.Username,
		Password: tools.NewGenPasswd(password),
		Nickname: r.RealName,
		Mail:     r.Email,
		Remark:   r.Remark,
		Status:   0,
	}

	err = isql.PendingUser.Add(pendingUser)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to submit registration: %v", err))
	}

	if err := tools.SendUserRegistrationPendingNotification(pendingUser.Username, pendingUser.Nickname, pendingUser.Mail); err != nil {
		common.Log.Warnf("Failed to send registration notification: %s, %v", pendingUser.Username, err)
	}

	return map[string]string{
		"message": "Registration submitted, pending admin approval",
	}, nil
}

// List returns pending user list
func (l PendingUserLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.PendingUserListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	pendingUsers, err := isql.PendingUser.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get pending users: %v", err))
	}

	count, err := isql.PendingUser.ListCount(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get pending user count: %v", err))
	}

	return response.PendingUserListRsp{
		Total:        int(count),
		PendingUsers: pendingUsers,
	}, nil
}

// Review approves or rejects pending user
func (l PendingUserLogic) Review(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.PendingUserReviewReq)
	if !ok {
		return nil, ReqAssertErr
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	pendingUser, err := isql.PendingUser.Find(r.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("pending user not found"))
	}

	now := time.Now()
	pendingUser.Status = r.Status
	pendingUser.Reviewer = ctxUser.Username
	pendingUser.ReviewRemark = r.ReviewRemark
	pendingUser.ReviewedAt = &now

	err = isql.PendingUser.Update(pendingUser)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update review status: %v", err))
	}

	// Approved: create user
	if r.Status == 1 {
		roleIds := r.RoleIds
		if len(roleIds) == 0 {
			roleIds = []uint{2}
		}
		roles, err := isql.Role.GetRolesByIds(roleIds)
		if err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("failed to get roles"))
		}

		user := &model.User{
			Username:      pendingUser.Username,
			Password:      pendingUser.Password,
			Nickname:      pendingUser.Nickname,
			GivenName:     pendingUser.Nickname,
			Mail:          pendingUser.Mail,
			Introduction:  pendingUser.Remark,
			Status:        1,
			Creator:       ctxUser.Username,
			Source:        "platform",
			Roles:         roles,
			UIDNumber:     r.UIDNumber,
			GIDNumber:     r.GIDNumber,
			HomeDirectory: r.HomeDirectory,
			LoginShell:    r.LoginShell,
			Mobile:        r.Mobile,
			JobNumber:     r.JobNumber,
			Position:      r.Position,
			PostalAddress: r.PostalAddress,
			Departments:   r.Departments,
		}
		
		if len(r.DepartmentId) > 0 {
			user.DepartmentId = tools.SliceToString(r.DepartmentId, ",")
		}

		// Apply pre-config if exists and request values are empty (pre-config as default)
		preConfig, _ := isql.UserPreConfig.GetByUsername(user.Username)
		if preConfig != nil {
			if user.Mail == "" && preConfig.Mail != "" {
				user.Mail = preConfig.Mail
			}
			if user.UIDNumber == 0 && preConfig.UIDNumber > 0 {
				user.UIDNumber = preConfig.UIDNumber
			}
			if user.GIDNumber == 0 && preConfig.GIDNumber > 0 {
				user.GIDNumber = preConfig.GIDNumber
			}
			if user.DepartmentId == "" && preConfig.DepartmentId != "" {
				user.DepartmentId = preConfig.DepartmentId
			}
			if user.Departments == "" && preConfig.Departments != "" {
				user.Departments = preConfig.Departments
			}
			if user.Mobile == "" && preConfig.Mobile != "" {
				user.Mobile = preConfig.Mobile
			}
			if user.JobNumber == "" && preConfig.JobNumber != "" {
				user.JobNumber = preConfig.JobNumber
			}
			if user.Position == "" && preConfig.Position != "" {
				user.Position = preConfig.Position
			}
			if user.PostalAddress == "" && preConfig.PostalAddress != "" {
				user.PostalAddress = preConfig.PostalAddress
			}
			if user.Introduction == "" && preConfig.Introduction != "" {
				user.Introduction = preConfig.Introduction
			}
			if user.HomeDirectory == "" && preConfig.HomeDirectory != "" {
				user.HomeDirectory = preConfig.HomeDirectory
			}
			if user.LoginShell == "" && preConfig.LoginShell != "" {
				user.LoginShell = preConfig.LoginShell
			}
			if user.Nickname == "" && preConfig.Nickname != "" {
				user.Nickname = preConfig.Nickname
			}
			if user.GivenName == "" && preConfig.GivenName != "" {
				user.GivenName = preConfig.GivenName
			}
			// Mark pre-config as used
			_ = isql.UserPreConfig.MarkAsUsed(preConfig.ID)
		}

		if user.GIDNumber == 0 && user.UIDNumber > 0 {
			user.GIDNumber = user.UIDNumber
		}
		if user.HomeDirectory == "" {
			user.HomeDirectory = fmt.Sprintf("/home/%s", user.Username)
		}
		if user.LoginShell == "" {
			user.LoginShell = "/bin/bash"
		}

		var groups []*model.Group
		if len(r.DepartmentId) > 0 {
			groups, err = isql.Group.GetGroupByIds(r.DepartmentId)
			if err != nil {
				return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %v", err))
			}
		}

		// Decrypt password for CommonAddUser to re-encrypt
		decryptedPassword := tools.NewParPasswd(pendingUser.Password)
		user.Password = decryptedPassword
		
		err = CommonAddUser(user, groups)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to create user: %v", err))
		}

		// Create group user permissions and sudo rules
		if len(groups) > 0 {
			for _, group := range groups {
				groupUserPerm := &model.GroupUserPermission{
					GroupID:     group.ID,
					UserID:      user.ID,
					AllowSudo:   r.AllowSudo,
					AllowSSHKey: r.AllowSSHKey,
					SudoRules:   r.SudoRules,
				}
				err = common.DB.Create(groupUserPerm).Error
				if err != nil {
					common.Log.Warnf("Failed to create group user permission: %v", err)
				}

				if r.AllowSudo {
					sudoRuleName := fmt.Sprintf("sudouser-%s-%s", user.Username, group.GroupName)
					sudoRule := &model.SudoRule{
						Name:        sudoRuleName,
						Description: fmt.Sprintf("Sudo permission for %s in %s", user.Username, group.GroupName),
						User:        user.Username,
						Host:        "ALL",
						Command:     "ALL",
						RunAsUser:   "ALL",
						RunAsGroup:  "ALL",
						Options:     func() string {
							if r.SudoRules != "" {
								return r.SudoRules
							}
							return "!authenticate"
						}(),
						Creator: ctxUser.Username,
					}
					
					err = ildap.Sudo.Add(sudoRule)
					if err != nil {
						common.Log.Warnf("Failed to create LDAP sudo rule: %v", err)
					} else {
						err = isql.SudoRule.Add(sudoRule)
						if err != nil {
							common.Log.Warnf("Failed to save sudo rule to MySQL: %v", err)
						}
					}
				}
			}
		}

		if err := tools.SendUserRegistrationApprovedNotification(user.Username, user.Nickname, user.Mail); err != nil {
			common.Log.Warnf("Failed to send approval notification: %s, %v", user.Username, err)
		}
	} else {
		// Rejected
		if err := tools.SendUserRegistrationRejectedNotification(pendingUser.Username, pendingUser.Nickname, pendingUser.Mail, r.ReviewRemark); err != nil {
			common.Log.Warnf("Failed to send rejection notification: %s, %v", pendingUser.Username, err)
		}
	}

	return map[string]string{
		"message": "Review completed",
	}, nil
}

// Delete removes pending users
func (l PendingUserLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.PendingUserDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.Ids {
		err := isql.PendingUser.Delete(id)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to delete pending user: %v", err))
		}
	}

	return nil, nil
}
