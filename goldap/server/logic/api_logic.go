package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type ApiLogic struct{}

// Add creates a new API record
func (l ApiLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	api := model.Api{
		Method:   r.Method,
		Path:     r.Path,
		Category: r.Category,
		Remark:   r.Remark,
		Creator:  ctxUser.Username,
	}

	if err := isql.Api.Add(&api); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to create API: %s", err.Error()))
	}

	return nil, nil
}

// List returns API list
func (l ApiLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	apis, err := isql.Api.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get API list: %s", err.Error()))
	}

	rets := make([]model.Api, 0)
	for _, api := range apis {
		rets = append(rets, *api)
	}

	count, err := isql.Api.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get API count"))
	}

	return response.ApiListRsp{
		Total: count,
		Apis:  rets,
	}, nil
}

// GetTree returns API tree grouped by category
func (l ApiLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiGetTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	_ = r

	apis, err := isql.Api.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get API list: %s", err.Error()))
	}

	var categoryList []string
	for _, api := range apis {
		categoryList = append(categoryList, api.Category)
	}
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*response.ApiTreeRsp, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &response.ApiTreeRsp{
			ID:       -i,
			Remark:   category,
			Category: category,
			Children: nil,
		}
		for _, api := range apis {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, api)
			}
		}
	}

	return apiTree, nil
}

// Update modifies an API record
func (l ApiLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Api.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("API not found"))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	oldData := new(model.Api)
	if err := isql.Api.Find(filter, oldData); err != nil {
		return nil, tools.NewMySqlError(err)
	}

	api := model.Api{
		Model:    oldData.Model,
		Method:   r.Method,
		Path:     r.Path,
		Category: r.Category,
		Remark:   r.Remark,
		Creator:  ctxUser.Username,
	}

	if err := isql.Api.Update(&api); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update API: %s", err.Error()))
	}

	return nil, nil
}

// Delete removes API records
func (l ApiLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.ApiIds {
		filter := tools.H{"id": int(id)}
		if !isql.Api.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("API not found"))
		}
	}

	if err := isql.Api.Delete(r.ApiIds); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete API: %s", err.Error()))
	}

	return nil, nil
}
