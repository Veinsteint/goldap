package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type RoleLogic struct{}

// Add creates a new role
func (l RoleLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.Role.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewValidatorError(fmt.Errorf("role name already exists"))
	}

	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user role level: %s", err.Error()))
	}

	if minSort != 1 {
		return nil, tools.NewValidatorError(fmt.Errorf("permission denied"))
	}

	if minSort >= r.Sort {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot create role with equal or higher level"))
	}

	role := model.Role{
		Name:    r.Name,
		Keyword: r.Keyword,
		Remark:  r.Remark,
		Status:  r.Status,
		Sort:    r.Sort,
		Creator: ctxUser.Username,
	}

	if err := isql.Role.Add(&role); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to create role: %s", err.Error()))
	}

	return nil, nil
}

// List returns role list
func (l RoleLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	roles, err := isql.Role.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get role list: %s", err.Error()))
	}

	count, err := isql.Role.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get role count"))
	}

	rets := make([]model.Role, 0)
	for _, role := range roles {
		rets = append(rets, *role)
	}

	return response.RoleListRsp{
		Total: count,
		Roles: rets,
	}, nil
}

// Update modifies a role
func (l RoleLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.ID}
	if !isql.Role.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("role not found"))
	}

	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user role level: %s", err.Error()))
	}

	if minSort != 1 {
		return nil, tools.NewValidatorError(fmt.Errorf("permission denied"))
	}

	roles, _ := isql.Role.GetRolesByIds([]uint{r.ID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlError(fmt.Errorf("role not found"))
	}

	if minSort >= roles[0].Sort {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot update role with equal or higher level"))
	}

	if minSort >= r.Sort {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot set role level equal or higher than self"))
	}

	oldData := new(model.Role)
	if err := isql.Role.Find(filter, oldData); err != nil {
		return nil, tools.NewMySqlError(err)
	}

	role := model.Role{
		Model:   oldData.Model,
		Name:    r.Name,
		Keyword: r.Keyword,
		Remark:  r.Remark,
		Status:  r.Status,
		Sort:    r.Sort,
		Creator: ctxUser.Username,
	}

	if err := isql.Role.Update(&role); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update role: %s", err.Error()))
	}

	// Update casbin policies if keyword changed
	if r.Keyword != roles[0].Keyword {
		rolePolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) > 0 {
			rolePoliciesCopy := make([][]string, 0)
			for _, policy := range rolePolicies {
				policyCopy := make([]string, len(policy))
				copy(policyCopy, policy)
				rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
				policy[0] = r.Keyword
			}

			if isAdded, _ := common.CasbinEnforcer.AddPolicies(rolePolicies); !isAdded {
				return nil, tools.NewOperationError(fmt.Errorf("failed to update role policies"))
			}
			if isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rolePoliciesCopy); !isRemoved {
				return nil, tools.NewOperationError(fmt.Errorf("failed to update role policies"))
			}
			if err := common.CasbinEnforcer.LoadPolicy(); err != nil {
				return nil, tools.NewOperationError(fmt.Errorf("failed to reload policies"))
			}
		}
	}

	isql.User.ClearUserInfoCache()
	return nil, nil
}

// Delete removes roles
func (l RoleLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user role level: %s", err.Error()))
	}

	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewMySqlError(fmt.Errorf("roles not found"))
	}

	for _, role := range roles {
		if minSort >= role.Sort {
			return nil, tools.NewValidatorError(fmt.Errorf("cannot delete role with equal or higher level"))
		}
	}

	if err := isql.Role.Delete(r.RoleIds); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete role: %s", err.Error()))
	}

	isql.User.ClearUserInfoCache()
	return nil, nil
}

// GetMenuList returns menus assigned to role
func (l RoleLogic) GetMenuList(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleGetMenuListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	menus, err := isql.Role.GetRoleMenusById(r.RoleID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get role menus: %s", err.Error()))
	}

	return menus, nil
}

// GetApiList returns APIs assigned to role
func (l RoleLogic) GetApiList(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleGetApiListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	role := new(model.Role)
	if err := isql.Role.Find(tools.H{"id": r.RoleID}, role); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get role: %s", err.Error()))
	}

	policies := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)

	apis, err := isql.Api.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get API list: %s", err.Error()))
	}

	accessApis := make([]*model.Api, 0)
	for _, policy := range policies {
		for _, api := range apis {
			if policy[1] == api.Path && policy[2] == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, nil
}

// UpdateMenus updates role menus
func (l RoleLogic) UpdateMenus(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateMenusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	roles, _ := isql.Role.GetRolesByIds([]uint{r.RoleID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlError(fmt.Errorf("role not found"))
	}

	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user role level: %s", err.Error()))
	}

	if minSort != 1 && minSort >= roles[0].Sort {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot update role with equal or higher level"))
	}

	ctxUserMenus, err := isql.Menu.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user menus: %s", err.Error()))
	}

	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	reqMenus := make([]*model.Menu, 0)

	if minSort != 1 {
		for _, id := range r.MenuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				return nil, tools.NewValidatorError(fmt.Errorf("no permission to set menu ID %d", id))
			}
		}
		for _, id := range r.MenuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		menus, err := isql.Menu.List()
		if err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("failed to get menu list: %s", err.Error()))
		}
		for _, menuId := range r.MenuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	if err := isql.Role.UpdateRoleMenus(roles[0]); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update role menus: %s", err.Error()))
	}

	return nil, nil
}

// UpdateApis updates role APIs
func (l RoleLogic) UpdateApis(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateApisReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	roles, _ := isql.Role.GetRolesByIds([]uint{r.RoleID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlError(fmt.Errorf("role not found"))
	}

	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get user role level: %s", err.Error()))
	}

	if minSort != 1 && minSort >= roles[0].Sort {
		return nil, tools.NewValidatorError(fmt.Errorf("cannot update role with equal or higher level"))
	}

	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	apis, err := isql.Api.GetApisById(r.ApiIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get APIs"))
	}

	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{roles[0].Keyword, api.Path, api.Method})
	}

	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				return nil, tools.NewValidatorError(fmt.Errorf("no permission to set API %s %s", reqPolicy[1], reqPolicy[2]))
			}
		}
	}

	if err := isql.Role.UpdateRoleApis(roles[0].Keyword, reqRolePolicies); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update role APIs"))
	}

	return nil, nil
}
