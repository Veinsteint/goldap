package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type MenuLogic struct{}

// Add creates a new menu
func (l MenuLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.Menu.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewMySqlError(fmt.Errorf("menu name already exists"))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	menu := model.Menu{
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	if err := isql.Menu.Add(&menu); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to create menu: %s", err.Error()))
	}

	return nil, nil
}

// Update modifies an existing menu
func (l MenuLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Menu.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("menu not found"))
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	oldData := new(model.Menu)
	if err := isql.Menu.Find(filter, oldData); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get menu: %s", err.Error()))
	}

	menu := model.Menu{
		Model:      oldData.Model,
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	if err := isql.Menu.Update(&menu); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update menu: %s", err.Error()))
	}

	return nil, nil
}

// Delete removes menu records
func (l MenuLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.MenuIds {
		filter := tools.H{"id": int(id)}
		if !isql.Menu.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("menu not found"))
		}
	}

	if err := isql.Menu.Delete(r.MenuIds); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete menu: %s", err.Error()))
	}

	return nil, nil
}

// GetTree returns menu tree
func (l MenuLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.MenuGetTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	menus, err := isql.Menu.List()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get menu list: %s", err.Error()))
	}

	return isql.GenMenuTree(0, menus), nil
}

// GetAccessTree returns user's accessible menu tree
func (l MenuLogic) GetAccessTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuGetAccessTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.ID}
	if !isql.User.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not found"))
	}

	user := new(model.User)
	if err := isql.User.Find(filter, user); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user: %s", err.Error()))
	}

	var roleIds []uint
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.ID)
	}

	menus, err := isql.Menu.ListUserMenus(roleIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get menu list: %s", err.Error()))
	}

	return isql.GenMenuTree(0, menus), nil
}
