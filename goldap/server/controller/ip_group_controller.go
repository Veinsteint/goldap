package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type IPGroupController struct{}

// Add creates IP group
// @Summary Create IP Group
// @Description Create new IP group
// @Tags IP Group Management
// @Accept json
// @Produce json
// @Param data body request.IPGroupAddReq true "IP group data"
// @Router /ip-group [post]
func (ic *IPGroupController) Add(c *gin.Context) {
	req := new(request.IPGroupAddReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroup.Add(c, req)
	})
}

// Update updates IP group
// @Summary Update IP Group
// @Description Update existing IP group
// @Tags IP Group Management
// @Accept json
// @Produce json
// @Param data body request.IPGroupUpdateReq true "IP group data"
// @Router /ip-group [put]
func (ic *IPGroupController) Update(c *gin.Context) {
	req := new(request.IPGroupUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroup.Update(c, req)
	})
}

// Delete removes IP group
// @Summary Delete IP Group
// @Description Delete IP group
// @Tags IP Group Management
// @Accept json
// @Produce json
// @Param data body request.IPGroupDeleteReq true "IP group ID"
// @Router /ip-group [delete]
func (ic *IPGroupController) Delete(c *gin.Context) {
	req := new(request.IPGroupDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.IPGroup.Delete(c, req)
	})
}

// List returns IP group list
// @Summary Get IP Group List
// @Description Get IP group list
// @Tags IP Group Management
// @Produce json
// @Router /ip-group [get]
func (ic *IPGroupController) List(c *gin.Context) {
	req := struct{}{}
	Run(c, req, func() (any, any) {
		return logic.IPGroup.List(c, req)
	})
}
