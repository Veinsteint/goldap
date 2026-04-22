package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type ApiController struct{}

// List returns API list
// @Summary Get API List
// @Description Get API endpoint list
// @Tags API Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /api/list [get]
// @Security ApiKeyAuth
func (m *ApiController) List(c *gin.Context) {
	req := new(request.ApiListReq)
	Run(c, req, func() (any, any) {
		return logic.Api.List(c, req)
	})
}

// GetTree returns API tree structure
// @Summary Get API Tree
// @Description Get API endpoint tree structure
// @Tags API Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /api/tree [get]
// @Security ApiKeyAuth
func (m *ApiController) GetTree(c *gin.Context) {
	req := new(request.ApiGetTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Api.GetTree(c, req)
	})
}

// Add creates new API endpoint
// @Summary Create API
// @Description Create new API endpoint
// @Tags API Management
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiAddReq true "API data"
// @Success 200 {object} response.ResponseBody
// @Router /api/add [post]
// @Security ApiKeyAuth
func (m *ApiController) Add(c *gin.Context) {
	req := new(request.ApiAddReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Add(c, req)
	})
}

// Update updates API endpoint
// @Summary Update API
// @Description Update existing API endpoint
// @Tags API Management
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiUpdateReq true "API data"
// @Success 200 {object} response.ResponseBody
// @Router /api/update [post]
// @Security ApiKeyAuth
func (m *ApiController) Update(c *gin.Context) {
	req := new(request.ApiUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Update(c, req)
	})
}

// Delete removes API endpoint
// @Summary Delete API
// @Description Delete API endpoint
// @Tags API Management
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiDeleteReq true "API ID"
// @Success 200 {object} response.ResponseBody
// @Router /api/delete [post]
// @Security ApiKeyAuth
func (m *ApiController) Delete(c *gin.Context) {
	req := new(request.ApiDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Delete(c, req)
	})
}
