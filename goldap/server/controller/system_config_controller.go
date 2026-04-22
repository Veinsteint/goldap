package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type SystemConfigController struct{}

// Get retrieves system configuration
// @Summary Get System Configuration
// @Description Get system configuration settings
// @Tags System Configuration
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /system/config/get [get]
// @Security ApiKeyAuth
func (m *SystemConfigController) Get(c *gin.Context) {
	req := new(request.SystemConfigGetReq)
	Run(c, req, func() (any, any) {
		return logic.SystemConfig.Get(c, req)
	})
}

// Update updates system configuration
// @Summary Update System Configuration
// @Description Update system configuration settings
// @Tags System Configuration
// @Accept application/json
// @Produce application/json
// @Param data body request.SystemConfigUpdateReq true "System configuration data"
// @Success 200 {object} response.ResponseBody
// @Router /system/config/update [post]
// @Security ApiKeyAuth
func (m *SystemConfigController) Update(c *gin.Context) {
	req := new(request.SystemConfigUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.SystemConfig.Update(c, req)
	})
}

