package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type FieldRelationLogic struct{}

// Add creates a field relation
func (l FieldRelationLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.FieldRelation.Exist(tools.H{"flag": r.Flag}) {
		return nil, tools.NewValidatorError(fmt.Errorf("field relation already exists"))
	}

	attr, err := tools.MapToJson(r.Attributes)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to convert map to JSON: %s", err.Error()))
	}

	frObj := model.FieldRelation{
		Flag:       r.Flag,
		Attributes: datatypes.JSON(attr),
	}

	if err := isql.FieldRelation.Add(&frObj); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to create field relation: %s", err.Error()))
	}

	return nil, nil
}

// List returns field relation list
func (l FieldRelationLogic) List(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.FieldRelationListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	frs, err := isql.FieldRelation.List()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get field relations: %s", err.Error()))
	}

	return frs, nil
}

// Update modifies a field relation
func (l FieldRelationLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"flag": r.Flag}
	if !isql.FieldRelation.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("field relation not found"))
	}

	oldData := new(model.FieldRelation)
	if err := isql.FieldRelation.Find(filter, oldData); err != nil {
		return nil, tools.NewMySqlError(err)
	}

	attr, err := tools.MapToJson(r.Attributes)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to convert map to JSON: %s", err.Error()))
	}

	frObj := model.FieldRelation{
		Model:      oldData.Model,
		Flag:       r.Flag,
		Attributes: datatypes.JSON(attr),
	}

	if err := isql.FieldRelation.Update(&frObj); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update field relation: %s", err.Error()))
	}

	return nil, nil
}

// Delete removes field relations
func (l FieldRelationLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.FieldRelationIds {
		filter := tools.H{"id": int(id)}
		if !isql.FieldRelation.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("field relation not found"))
		}
	}

	if err := isql.FieldRelation.Delete(r.FieldRelationIds); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete field relation: %s", err.Error()))
	}

	return nil, nil
}
