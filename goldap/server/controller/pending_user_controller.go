package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type PendingUserController struct{}

// Register handles user registration
// @Summary User Registration
// @Description User registration request (requires admin approval)
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param data body request.UserRegisterReq true "Registration data"
// @Success 200 {object} response.ResponseBody
// @Router /base/register [post]
func (m *PendingUserController) Register(c *gin.Context) {
	req := new(request.UserRegisterReq)
	Run(c, req, func() (any, any) {
		return logic.PendingUser.Register(c, req)
	})
}

// List returns pending user list
// @Summary Get Pending User List
// @Description Get list of users pending approval
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data query request.PendingUserListReq false "Query parameters"
// @Success 200 {object} response.ResponseBody
// @Router /user/pending/list [get]
// @Security ApiKeyAuth
func (m *PendingUserController) List(c *gin.Context) {
	req := new(request.PendingUserListReq)
	Run(c, req, func() (any, any) {
		return logic.PendingUser.List(c, req)
	})
}

// Review approves or rejects pending user
// @Summary Review Pending User
// @Description Approve or reject pending user registration
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.PendingUserReviewReq true "Review data"
// @Success 200 {object} response.ResponseBody
// @Router /user/pending/review [post]
// @Security ApiKeyAuth
func (m *PendingUserController) Review(c *gin.Context) {
	req := new(request.PendingUserReviewReq)
	Run(c, req, func() (any, any) {
		return logic.PendingUser.Review(c, req)
	})
}

// Delete removes pending users
// @Summary Delete Pending Users
// @Description Batch delete pending user registrations
// @Tags User Management
// @Accept application/json
// @Produce application/json
// @Param data body request.PendingUserDeleteReq true "Delete data"
// @Success 200 {object} response.ResponseBody
// @Router /user/pending/delete [post]
// @Security ApiKeyAuth
func (m *PendingUserController) Delete(c *gin.Context) {
	req := new(request.PendingUserDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.PendingUser.Delete(c, req)
	})
}
