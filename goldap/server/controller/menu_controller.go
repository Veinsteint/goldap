package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

// GetTree returns menu tree
// @Summary Get Menu Tree
// @Description Get menu tree structure
// @Tags Menu Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /menu/tree [get]
// @Security ApiKeyAuth
func (m *MenuController) GetTree(c *gin.Context) {
	req := new(request.MenuGetTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.GetTree(c, req)
	})
}

// GetAccessTree returns user's accessible menu tree
// @Summary Get User Menu Tree
// @Description Get menus accessible to current user
// @Tags Menu Management
// @Accept application/json
// @Produce application/json
// @Param id query int true "User ID"
// @Success 200 {object} response.ResponseBody
// @Router /menu/access/tree [get]
// @Security ApiKeyAuth
func (m *MenuController) GetAccessTree(c *gin.Context) {
	req := new(request.MenuGetAccessTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.GetAccessTree(c, req)
	})
}

// Add creates new menu
// @Summary Create Menu
// @Description Create new menu item
// @Tags Menu Management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuAddReq true "Menu data"
// @Success 200 {object} response.ResponseBody
// @Router /menu/add [post]
// @Security ApiKeyAuth
func (m *MenuController) Add(c *gin.Context) {
	req := new(request.MenuAddReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Add(c, req)
	})
}

// Update updates menu
// @Summary Update Menu
// @Description Update existing menu item
// @Tags Menu Management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuUpdateReq true "Menu data"
// @Success 200 {object} response.ResponseBody
// @Router /menu/update [post]
// @Security ApiKeyAuth
func (m *MenuController) Update(c *gin.Context) {
	req := new(request.MenuUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Update(c, req)
	})
}

// Delete removes menu
// @Summary Delete Menu
// @Description Delete menu item
// @Tags Menu Management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuDeleteReq true "Menu ID"
// @Success 200 {object} response.ResponseBody
// @Router /menu/delete [post]
// @Security ApiKeyAuth
func (m *MenuController) Delete(c *gin.Context) {
	req := new(request.MenuDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Delete(c, req)
	})
}
