package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// List returns role list
// @Summary Get Role List
// @Description Get role list
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /role/list [get]
// @Security ApiKeyAuth
func (m *RoleController) List(c *gin.Context) {
	req := new(request.RoleListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.List(c, req)
	})
}

// Add creates new role
// @Summary Create Role
// @Description Create new role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param data body request.RoleAddReq true "Role data"
// @Success 200 {object} response.ResponseBody
// @Router /role/add [post]
// @Security ApiKeyAuth
func (m *RoleController) Add(c *gin.Context) {
	req := new(request.RoleAddReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Add(c, req)
	})
}

// Update updates role
// @Summary Update Role
// @Description Update existing role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param data body request.RoleUpdateReq true "Role data"
// @Success 200 {object} response.ResponseBody
// @Router /role/update [post]
// @Security ApiKeyAuth
func (m *RoleController) Update(c *gin.Context) {
	req := new(request.RoleUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Update(c, req)
	})
}

// Delete removes role
// @Summary Delete Role
// @Description Delete role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param data body request.RoleDeleteReq true "Role ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/delete [post]
// @Security ApiKeyAuth
func (m *RoleController) Delete(c *gin.Context) {
	req := new(request.RoleDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Delete(c, req)
	})
}

// GetMenuList returns menus assigned to role
// @Summary Get Role Menus
// @Description Get menus assigned to role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "Role ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getmenulist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetMenuList(c *gin.Context) {
	req := new(request.RoleGetMenuListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetMenuList(c, req)
	})
}

// GetApiList returns APIs assigned to role
// @Summary Get Role APIs
// @Description Get APIs assigned to role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "Role ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getapilist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetApiList(c *gin.Context) {
	req := new(request.RoleGetApiListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetApiList(c, req)
	})
}

// UpdateMenus updates role menus
// @Summary Update Role Menus
// @Description Update menus assigned to role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param data body request.RoleUpdateMenusReq true "Menu IDs"
// @Success 200 {object} response.ResponseBody
// @Router /role/updatemenus [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateMenus(c *gin.Context) {
	req := new(request.RoleUpdateMenusReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateMenus(c, req)
	})
}

// UpdateApis updates role APIs
// @Summary Update Role APIs
// @Description Update APIs assigned to role
// @Tags Role Management
// @Accept application/json
// @Produce application/json
// @Param data body request.RoleUpdateApisReq true "API IDs"
// @Success 200 {object} response.ResponseBody
// @Router /role/updateapis [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateApis(c *gin.Context) {
	req := new(request.RoleUpdateApisReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateApis(c, req)
	})
}
