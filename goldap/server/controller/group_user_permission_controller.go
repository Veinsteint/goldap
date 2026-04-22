package controller

import (
	"fmt"

	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type GroupUserPermissionController struct{}

// Add creates group user permission
// @Summary Create Group User Permission
// @Description Create new group user permission
// @Tags Group User Permission
// @Accept json
// @Produce json
// @Param data body request.GroupUserPermissionAddReq true "Permission data"
// @Router /group-user-permission [post]
func (gc *GroupUserPermissionController) Add(c *gin.Context) {
	req := new(request.GroupUserPermissionAddReq)
	Run(c, req, func() (any, any) {
		return logic.GroupUserPermission.Add(c, req)
	})
}

// Update updates group user permission
// @Summary Update Group User Permission
// @Description Update existing group user permission
// @Tags Group User Permission
// @Accept json
// @Produce json
// @Param data body request.GroupUserPermissionUpdateReq true "Permission data"
// @Router /group-user-permission [put]
func (gc *GroupUserPermissionController) Update(c *gin.Context) {
	req := new(request.GroupUserPermissionUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.GroupUserPermission.Update(c, req)
	})
}

// Delete removes group user permission
// @Summary Delete Group User Permission
// @Description Delete group user permission
// @Tags Group User Permission
// @Accept json
// @Produce json
// @Param data body request.GroupUserPermissionDeleteReq true "Permission ID"
// @Router /group-user-permission [delete]
func (gc *GroupUserPermissionController) Delete(c *gin.Context) {
	req := new(request.GroupUserPermissionDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.GroupUserPermission.Delete(c, req)
	})
}

// List returns group user permission list
// @Summary Get Group User Permission List
// @Description Get group user permission list with filters
// @Tags Group User Permission
// @Produce json
// @Param groupId query uint false "Group ID"
// @Param userId query uint false "User ID"
// @Router /group-user-permission [get]
func (gc *GroupUserPermissionController) List(c *gin.Context) {
	req := &request.GroupUserPermissionListReq{}

	if groupIdStr := c.Query("groupId"); groupIdStr != "" {
		var groupId uint
		if _, err := fmt.Sscanf(groupIdStr, "%d", &groupId); err == nil {
			req.GroupID = groupId
		}
	}
	if userIdStr := c.Query("userId"); userIdStr != "" {
		var userId uint
		if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err == nil {
			req.UserID = userId
		}
	}

	Run(c, req, func() (any, any) {
		return logic.GroupUserPermission.List(c, req)
	})
}
