package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// Add creates a new user record
// @Summary Create user record
// @Description Create a new user record
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserAddReq true "User creation request"
// @Success 200 {object} response.ResponseBody
// @Router /user/add [post]
// @Security ApiKeyAuth
func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Run(c, req, func() (any, any) {
		return logic.User.Add(c, req)
	})
}

// Update modifies an existing user record
// @Summary Update user record
// @Description Update an existing user record
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserUpdateReq true "User update request"
// @Success 200 {object} response.ResponseBody
// @Router /user/update [post]
// @Security ApiKeyAuth
func (m *UserController) Update(c *gin.Context) {
	req := new(request.UserUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.User.Update(c, req)
	})
}

// List retrieves all user records
// @Summary Get all user records
// @Description Get list of all user records
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/list [get]
// @Security ApiKeyAuth
func (m *UserController) List(c *gin.Context) {
	req := new(request.UserListReq)
	Run(c, req, func() (any, any) {
		return logic.User.List(c, req)
	})
}

// Delete removes a user record
// @Summary Delete user record
// @Description Delete a user record
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserDeleteReq true "User deletion request"
// @Success 200 {object} response.ResponseBody
// @Router /user/delete [post]
// @Security ApiKeyAuth
func (m UserController) Delete(c *gin.Context) {
	req := new(request.UserDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.User.Delete(c, req)
	})
}

// ChangePwd updates user password
// @Summary Update password
// @Description Update user password
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserChangePwdReq true "Password change request"
// @Success 200 {object} response.ResponseBody
// @Router /user/changePwd [post]
// @Security ApiKeyAuth
func (m UserController) ChangePwd(c *gin.Context) {
	req := new(request.UserChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangePwd(c, req)
	})
}

// ResetPassword resets user password
// @Summary Reset user password
// @Description Reset user password to random value and send email notification
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserResetPasswordReq true "Password reset request"
// @Success 200 {object} response.ResponseBody
// @Router /user/resetPassword [post]
// @Security ApiKeyAuth
func (m UserController) ResetPassword(c *gin.Context) {
	req := new(request.UserResetPasswordReq)
	Run(c, req, func() (any, any) {
		return logic.User.ResetPassword(c, req)
	})
}

// ChangeUserStatus changes user status
// @Summary Change user status
// @Description Change user status (enable/disable)
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.UserChangeUserStatusReq true "Status change request"
// @Success 200 {object} response.ResponseBody
// @Router /user/changeUserStatus [post]
// @Security ApiKeyAuth
func (m UserController) ChangeUserStatus(c *gin.Context) {
	req := new(request.UserChangeUserStatusReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangeUserStatus(c, req)
	})
}

// GetUserInfo gets current logged-in user information
// @Summary Get current user info
// @Description Get information of currently logged-in user
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/info [get]
// @Security ApiKeyAuth
func (uc UserController) GetUserInfo(c *gin.Context) {
	req := new(request.UserGetUserInfoReq)
	Run(c, req, func() (any, any) {
		return logic.User.GetUserInfo(c, req)
	})
}

// SyncOpenLdapUsers syncs LDAP user information
// @Summary Sync LDAP users
// @Description Sync user information from LDAP
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.SyncOpenLdapUserReq true "LDAP sync request"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncOpenLdapUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncOpenLdapUsers(c *gin.Context) {
	req := new(request.SyncOpenLdapUserReq)
	Run(c, req, func() (any, any) {
		return logic.OpenLdap.SyncOpenLdapUsers(c, req)
	})
}

// SyncSqlUsers syncs SQL users to LDAP
// @Summary Sync SQL users to LDAP
// @Description Sync user information from SQL to LDAP
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.SyncSqlUserReq true "SQL sync request"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncSqlUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncSqlUsers(c *gin.Context) {
	req := new(request.SyncSqlUserReq)
	Run(c, req, func() (any, any) {
		return logic.Sql.SyncSqlUsers(c, req)
	})
}
