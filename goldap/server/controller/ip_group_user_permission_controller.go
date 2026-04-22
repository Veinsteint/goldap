package controller

import (
	"fmt"

	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type IPGroupUserPermissionController struct{}

// Add creates IP group user permission
// @Summary Create IP Group User Permission
// @Description Create new IP group user permission
// @Tags IP Group User Permission
// @Accept json
// @Produce json
// @Param data body request.IPGroupUserPermissionAddReq true "Permission data"
// @Router /ip-group-user-permission [post]
func (ic *IPGroupUserPermissionController) Add(c *gin.Context) {
	req := new(request.IPGroupUserPermissionAddReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroupUserPermission.Add(c, req)
	})
}

// Update updates IP group user permission
// @Summary Update IP Group User Permission
// @Description Update existing IP group user permission
// @Tags IP Group User Permission
// @Accept json
// @Produce json
// @Param data body request.IPGroupUserPermissionUpdateReq true "Permission data"
// @Router /ip-group-user-permission [put]
func (ic *IPGroupUserPermissionController) Update(c *gin.Context) {
	req := new(request.IPGroupUserPermissionUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroupUserPermission.Update(c, req)
	})
}

// Delete removes IP group user permission
// @Summary Delete IP Group User Permission
// @Description Delete IP group user permission
// @Tags IP Group User Permission
// @Accept json
// @Produce json
// @Param data body request.IPGroupUserPermissionDeleteReq true "Permission ID"
// @Router /ip-group-user-permission [delete]
func (ic *IPGroupUserPermissionController) Delete(c *gin.Context) {
	req := new(request.IPGroupUserPermissionDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroupUserPermission.Delete(c, req)
	})
}

// List returns IP group user permission list
// @Summary Get IP Group User Permission List
// @Description Get IP group user permission list with filters
// @Tags IP Group User Permission
// @Produce json
// @Param ipGroupId query uint false "IP Group ID"
// @Param userId query uint false "User ID"
// @Router /ip-group-user-permission [get]
func (ic *IPGroupUserPermissionController) List(c *gin.Context) {
	req := &request.IPGroupUserPermissionListReq{}

	if ipGroupIdStr := c.Query("ipGroupId"); ipGroupIdStr != "" {
		var ipGroupId uint
		if _, err := fmt.Sscanf(ipGroupIdStr, "%d", &ipGroupId); err == nil {
			req.IPGroupID = ipGroupId
		}
	}
	if userIdStr := c.Query("userId"); userIdStr != "" {
		var userId uint
		if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err == nil {
			req.UserID = userId
		}
	}

	Run(c, req, func() (any, any) {
		return logic.IPGroupUserPermission.List(c, req)
	})
}
