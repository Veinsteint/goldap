package logic

import (
	"fmt"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type UserPreConfigLogic struct{}

// List returns paginated user pre-config list
func (l UserPreConfigLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserPreConfigListReq)
	if !ok {
		return nil, ReqAssertErr
	}

	list, err := isql.UserPreConfig.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get list: %v", err))
	}

	count, err := isql.UserPreConfig.Count(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get count: %v", err))
	}

	return tools.H{"list": list, "total": count}, nil
}

// Add creates a new user pre-config
func (l UserPreConfigLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserPreConfigAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	// Check if username already exists
	if isql.UserPreConfig.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("username [%s] already has pre-config", r.Username))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user: %v", err))
	}

	preConfig := model.UserPreConfig{
		Username:      r.Username,
		Mail:          r.Mail,
		UIDNumber:     r.UIDNumber,
		GIDNumber:     r.GIDNumber,
		DepartmentId:  r.DepartmentId,
		Departments:   r.Departments,
		Mobile:        r.Mobile,
		JobNumber:     r.JobNumber,
		Position:      r.Position,
		PostalAddress: r.PostalAddress,
		Introduction:  r.Introduction,
		HomeDirectory: r.HomeDirectory,
		LoginShell:    r.LoginShell,
		Nickname:      r.Nickname,
		GivenName:     r.GivenName,
		Remark:        r.Remark,
		Creator:       ctxUser.Username,
	}

	// Set defaults
	if preConfig.HomeDirectory == "" && preConfig.Username != "" {
		preConfig.HomeDirectory = fmt.Sprintf("/home/%s", preConfig.Username)
	}
	if preConfig.LoginShell == "" {
		preConfig.LoginShell = "/bin/bash"
	}

	if err := isql.UserPreConfig.Add(&preConfig); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to create pre-config: %v", err))
	}

	return nil, nil
}

// Update modifies an existing user pre-config
func (l UserPreConfigLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserPreConfigUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if !isql.UserPreConfig.Exist(tools.H{"id": r.ID}) {
		return nil, tools.NewMySqlError(fmt.Errorf("pre-config not found"))
	}

	// Check for duplicate username (excluding current record)
	var existing model.UserPreConfig
	if err := isql.UserPreConfig.Find(tools.H{"username": r.Username}, &existing); err == nil {
		if existing.ID != r.ID {
			return nil, tools.NewValidatorError(fmt.Errorf("username [%s] already has pre-config", r.Username))
		}
	}

	preConfig := model.UserPreConfig{
		Username:      r.Username,
		Mail:          r.Mail,
		UIDNumber:     r.UIDNumber,
		GIDNumber:     r.GIDNumber,
		DepartmentId:  r.DepartmentId,
		Departments:   r.Departments,
		Mobile:        r.Mobile,
		JobNumber:     r.JobNumber,
		Position:      r.Position,
		PostalAddress: r.PostalAddress,
		Introduction:  r.Introduction,
		HomeDirectory: r.HomeDirectory,
		LoginShell:    r.LoginShell,
		Nickname:      r.Nickname,
		GivenName:     r.GivenName,
		Remark:        r.Remark,
	}
	preConfig.ID = r.ID

	if err := isql.UserPreConfig.Update(&preConfig); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update pre-config: %v", err))
	}

	// Sync to existing user if exists
	var user model.User
	if err := isql.User.Find(tools.H{"username": r.Username}, &user); err == nil {
		// User exists, sync the pre-config data to user
		syncUserFromPreConfig(&user, &preConfig)
		_ = isql.User.Update(&user, true) // Skip LDAP sync here, will be done separately
	}

	return nil, nil
}

// Delete removes user pre-configs
func (l UserPreConfigLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserPreConfigDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if err := isql.UserPreConfig.Delete(r.Ids); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete pre-configs: %v", err))
	}

	return nil, nil
}

// GetByUsername retrieves pre-config by username
func (l UserPreConfigLogic) GetByUsername(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserPreConfigGetByUsernameReq)
	if !ok {
		return nil, ReqAssertErr
	}

	config, err := isql.UserPreConfig.GetByUsername(r.Username)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get pre-config: %v", err))
	}

	return config, nil
}

// syncUserFromPreConfig updates user fields from pre-config (helper function)
func syncUserFromPreConfig(user *model.User, config *model.UserPreConfig) {
	if config.Mail != "" {
		user.Mail = config.Mail
	}
	if config.UIDNumber > 0 {
		user.UIDNumber = config.UIDNumber
	}
	if config.GIDNumber > 0 {
		user.GIDNumber = config.GIDNumber
	}
	if config.DepartmentId != "" {
		user.DepartmentId = config.DepartmentId
	}
	if config.Departments != "" {
		user.Departments = config.Departments
	}
	if config.Mobile != "" {
		user.Mobile = config.Mobile
	}
	if config.JobNumber != "" {
		user.JobNumber = config.JobNumber
	}
	if config.Position != "" {
		user.Position = config.Position
	}
	if config.PostalAddress != "" {
		user.PostalAddress = config.PostalAddress
	}
	if config.Introduction != "" {
		user.Introduction = config.Introduction
	}
	if config.HomeDirectory != "" {
		user.HomeDirectory = config.HomeDirectory
	}
	if config.LoginShell != "" {
		user.LoginShell = config.LoginShell
	}
	if config.Nickname != "" {
		user.Nickname = config.Nickname
	}
	if config.GivenName != "" {
		user.GivenName = config.GivenName
	}
}

// GetValidUsernames returns all pre-configured usernames with their registration status
// Status: "已注册-有效", "已注册-失效", "待审核", "未注册"
func (l UserPreConfigLogic) GetValidUsernames(c *gin.Context, req any) (data any, rspError any) {
	_ = c
	_ = req

	// Get all pre-configs without pagination limit
	preConfigs, err := isql.UserPreConfig.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get pre-configs: %v", err))
	}

	type UsernameStatus struct {
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Status   string `json:"status"`
		Remark   string `json:"remark"`
	}

	result := make([]UsernameStatus, 0, len(preConfigs))

	for _, pc := range preConfigs {
		status := "未注册"
		
		// First check if is_used is marked
		if pc.IsUsed {
			var user model.User
			if err := isql.User.Find(tools.H{"username": pc.Username}, &user); err == nil {
				if user.Status == 1 {
					status = "已注册-有效"
				} else {
					status = "已注册-失效"
				}
			} else {
				// is_used=1 but user doesn't exist
				if isql.PendingUser.Exist(tools.H{"username": pc.Username}) {
					status = "待审核"
				} else {
					status = "已使用(禁止注册)"
				}
			}
		} else {
			var user model.User
			if err := isql.User.Find(tools.H{"username": pc.Username}, &user); err == nil {
				// User exists
				if user.Status == 1 {
					status = "已注册-有效"
				} else {
					status = "已注册-失效"
				}
			} else {
				if isql.PendingUser.Exist(tools.H{"username": pc.Username}) {
					status = "待审核"
				}
				// Otherwise remains "未注册"
			}
		}

		result = append(result, UsernameStatus{
			Username: pc.Username,
			Nickname: pc.Nickname,
			Status:   status,
			Remark:   pc.Remark,
		})
	}

	return result, nil
}

// GetRegistrationMode returns the current registration mode
func (l UserPreConfigLogic) GetRegistrationMode(c *gin.Context, req any) (data any, rspError any) {
	_ = c
	_ = req

	mode := "preconfig" // default
	if config.Conf.Registration != nil && config.Conf.Registration.Mode != "" {
		mode = config.Conf.Registration.Mode
	}

	return map[string]string{"mode": mode}, nil
}

// SyncExistingUsers syncs all existing users (except admin) to pre-config table
func (l UserPreConfigLogic) SyncExistingUsers(c *gin.Context, req any) (data any, rspError any) {
	_ = c
	_ = req

	// Get all users except admin
	var users []*model.User
	userListReq := &request.UserListReq{PageNum: 1, PageSize: 10000}
	users, err := isql.User.List(userListReq)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %v", err))
	}

	syncedCount := 0
	skippedCount := 0

	for _, user := range users {
		// Skip admin user
		if user.Username == "admin" {
			skippedCount++
			continue
		}

		// Check if pre-config already exists
		existing, _ := isql.UserPreConfig.GetByUsername(user.Username)
		if existing != nil {
			skippedCount++
			continue
		}

		// Create new pre-config from user data
		preConfig := &model.UserPreConfig{
			Username:      user.Username,
			Mail:          user.Mail,
			UIDNumber:     user.UIDNumber,
			GIDNumber:     user.GIDNumber,
			DepartmentId:  user.DepartmentId,
			Departments:   user.Departments,
			Mobile:        user.Mobile,
			JobNumber:     user.JobNumber,
			Position:      user.Position,
			PostalAddress: user.PostalAddress,
			Introduction:  user.Introduction,
			HomeDirectory: user.HomeDirectory,
			LoginShell:    user.LoginShell,
			Nickname:      user.Nickname,
			GivenName:     user.GivenName,
			Creator:       "system",
			IsUsed:        true,
		}

		if err := isql.UserPreConfig.Add(preConfig); err != nil {
			continue
		}
		syncedCount++
	}

	return map[string]int{
		"synced":  syncedCount,
		"skipped": skippedCount,
	}, nil
}

// ApplyPreConfigToUser applies pre-config values to a user (exported for use in pending_user_logic)
func ApplyPreConfigToUser(user *model.User) bool {
	config, err := isql.UserPreConfig.GetByUsername(user.Username)
	if err != nil || config == nil {
		return false
	}

	syncUserFromPreConfig(user, config)
	return true
}

