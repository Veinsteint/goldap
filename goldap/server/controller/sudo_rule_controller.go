package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type SudoRuleController struct{}

// Add creates sudo rule
// @Summary Create Sudo Rule
// @Description Create new sudo rule
// @Tags Sudo Rule Management
// @Accept json
// @Produce json
// @Param data body request.SudoRuleAddReq true "Sudo rule data"
// @Router /sudo-rule [post]
func (sc *SudoRuleController) Add(c *gin.Context) {
	req := new(request.SudoRuleAddReq)
	Run(c, req, func() (any, any) {
		return logic.SudoRule.Add(c, req)
	})
}

// Update updates sudo rule
// @Summary Update Sudo Rule
// @Description Update existing sudo rule
// @Tags Sudo Rule Management
// @Accept json
// @Produce json
// @Param data body request.SudoRuleUpdateReq true "Sudo rule data"
// @Router /sudo-rule [put]
func (sc *SudoRuleController) Update(c *gin.Context) {
	req := new(request.SudoRuleUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.SudoRule.Update(c, req)
	})
}

// Delete removes sudo rule
// @Summary Delete Sudo Rule
// @Description Delete sudo rule
// @Tags Sudo Rule Management
// @Accept json
// @Produce json
// @Param data body request.SudoRuleDeleteReq true "Sudo rule ID"
// @Router /sudo-rule [delete]
func (sc *SudoRuleController) Delete(c *gin.Context) {
	req := new(request.SudoRuleDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.SudoRule.Delete(c, req)
	})
}

// List returns sudo rule list
// @Summary Get Sudo Rule List
// @Description Get sudo rule list
// @Tags Sudo Rule Management
// @Produce json
// @Router /sudo-rule [get]
func (sc *SudoRuleController) List(c *gin.Context) {
	req := struct{}{}
	Run(c, req, func() (any, any) {
		return logic.SudoRule.List(c, req)
	})
}
