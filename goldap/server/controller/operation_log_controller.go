package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type OperationLogController struct{}

// List returns operation log list
// @Summary Get Operation Log List
// @Description Get operation log list with filters
// @Tags Operation Log
// @Accept application/json
// @Produce application/json
// @Param username query string false "Username"
// @Param ip query string false "IP Address"
// @Param path query string false "Path"
// @Param method query string false "HTTP Method"
// @Param status query int false "Status Code"
// @Param pageNum query int false "Page Number"
// @Param pageSize query int false "Page Size"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/list [get]
// @Security ApiKeyAuth
func (m *OperationLogController) List(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.List(c, req)
	})
}

// Delete removes operation log entries
// @Summary Delete Operation Log
// @Description Delete operation log entries by IDs
// @Tags Operation Log
// @Accept application/json
// @Produce application/json
// @Param data body request.OperationLogDeleteReq true "Log IDs"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/delete [post]
// @Security ApiKeyAuth
func (m *OperationLogController) Delete(c *gin.Context) {
	req := new(request.OperationLogDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.Delete(c, req)
	})
}

// Clean clears all operation logs
// @Summary Clear Operation Logs
// @Description Clear all operation log entries
// @Tags Operation Log
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/clean [delete]
// @Security ApiKeyAuth
func (m *OperationLogController) Clean(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.Clean(c, req)
	})
}
