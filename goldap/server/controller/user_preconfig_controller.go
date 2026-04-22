package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type UserPreConfigController struct{}

// List returns user pre-config list
// @Summary Get User Pre-config List
// @Description Get user pre-config list with pagination
// @Tags User Pre-config
// @Accept application/json
// @Produce application/json
// @Param username query string false "Username filter"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/list [get]
// @Security ApiKeyAuth
func (m *UserPreConfigController) List(c *gin.Context) {
	req := new(request.UserPreConfigListReq)
	Run(c, req, func() (any, any) {
		return logic.UserPreConfig.List(c, req)
	})
}

// Add creates user pre-config
// @Summary Create User Pre-config
// @Description Create new user pre-config
// @Tags User Pre-config
// @Accept application/json
// @Produce application/json
// @Param data body request.UserPreConfigAddReq true "Pre-config data"
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/add [post]
// @Security ApiKeyAuth
func (m *UserPreConfigController) Add(c *gin.Context) {
	req := new(request.UserPreConfigAddReq)
	Run(c, req, func() (any, any) {
		return logic.UserPreConfig.Add(c, req)
	})
}

// Update updates user pre-config
// @Summary Update User Pre-config
// @Description Update existing user pre-config
// @Tags User Pre-config
// @Accept application/json
// @Produce application/json
// @Param data body request.UserPreConfigUpdateReq true "Pre-config data"
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/update [post]
// @Security ApiKeyAuth
func (m *UserPreConfigController) Update(c *gin.Context) {
	req := new(request.UserPreConfigUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.UserPreConfig.Update(c, req)
	})
}

// Delete removes user pre-config
// @Summary Delete User Pre-config
// @Description Delete user pre-config
// @Tags User Pre-config
// @Accept application/json
// @Produce application/json
// @Param data body request.UserPreConfigDeleteReq true "Pre-config IDs"
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/delete [post]
// @Security ApiKeyAuth
func (m *UserPreConfigController) Delete(c *gin.Context) {
	req := new(request.UserPreConfigDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.UserPreConfig.Delete(c, req)
	})
}

// GetByUsername returns pre-config by username
// @Summary Get User Pre-config by Username
// @Description Get user pre-config by username
// @Tags User Pre-config
// @Accept application/json
// @Produce application/json
// @Param username query string true "Username"
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/getByUsername [get]
// @Security ApiKeyAuth
func (m *UserPreConfigController) GetByUsername(c *gin.Context) {
	req := new(request.UserPreConfigGetByUsernameReq)
	Run(c, req, func() (any, any) {
		return logic.UserPreConfig.GetByUsername(c, req)
	})
}

// GetValidUsernames returns all pre-configured usernames with their status (public API)
// @Summary Get Valid Usernames for Registration
// @Description Get all pre-configured usernames with their registration status
// @Tags Registration
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/validUsernames [get]
func (m *UserPreConfigController) GetValidUsernames(c *gin.Context) {
	Run(c, nil, func() (any, any) {
		return logic.UserPreConfig.GetValidUsernames(c, nil)
	})
}

// GetRegistrationMode returns the current registration mode (public API)
// @Summary Get Registration Mode
// @Description Get the current registration mode (open or preconfig)
// @Tags Registration
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/registrationMode [get]
func (m *UserPreConfigController) GetRegistrationMode(c *gin.Context) {
	Run(c, nil, func() (any, any) {
		return logic.UserPreConfig.GetRegistrationMode(c, nil)
	})
}

// SyncExistingUsers syncs existing users to pre-config table
// @Summary Sync Existing Users to Pre-config
// @Description Sync all existing users (except admin) to pre-config table
// @Tags User Pre-config
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/preconfig/syncUsers [post]
// @Security ApiKeyAuth
func (m *UserPreConfigController) SyncExistingUsers(c *gin.Context) {
	Run(c, nil, func() (any, any) {
		return logic.UserPreConfig.SyncExistingUsers(c, nil)
	})
}

