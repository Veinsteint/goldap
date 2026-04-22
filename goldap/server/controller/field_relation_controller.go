package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type FieldRelationController struct{}

// List returns field relation list
// @Summary Get Field Relation List
// @Description Get field relation list
// @Tags Field Relation
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/list [get]
// @Security ApiKeyAuth
func (m *FieldRelationController) List(c *gin.Context) {
	req := new(request.FieldRelationListReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.List(c, req)
	})
}

// Add creates field relation
// @Summary Create Field Relation
// @Description Create new field relation
// @Tags Field Relation
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationAddReq true "Field relation data"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/add [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Add(c *gin.Context) {
	req := new(request.FieldRelationAddReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Add(c, req)
	})
}

// Update updates field relation
// @Summary Update Field Relation
// @Description Update existing field relation
// @Tags Field Relation
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationUpdateReq true "Field relation data"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/update [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Update(c *gin.Context) {
	req := new(request.FieldRelationUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Update(c, req)
	})
}

// Delete removes field relation
// @Summary Delete Field Relation
// @Description Delete field relation
// @Tags Field Relation
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationDeleteReq true "Field relation ID"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/delete [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Delete(c *gin.Context) {
	req := new(request.FieldRelationDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Delete(c, req)
	})
}
