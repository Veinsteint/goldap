package isql

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"

	"gorm.io/gorm"
)

type RoleService struct{}

// Exist checks if role exists
func (s RoleService) Exist(filter map[string]any) bool {
	var dataObj model.Role
	err := common.DB.Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// List returns paginated role list
func (s RoleService) List(req *request.RoleListReq) ([]*model.Role, error) {
	var list []*model.Role
	db := common.DB.Model(&model.Role{}).Order("created_at DESC")

	if name := strings.TrimSpace(req.Name); name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// Count returns total role count
func (s RoleService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Role{}).Count(&count).Error
	return count, err
}

// Add creates a new role
func (s RoleService) Add(role *model.Role) error {
	return common.DB.Create(role).Error
}

// Update modifies role
func (s RoleService) Update(role *model.Role) error {
	return common.DB.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// Find gets single role by filter
func (s RoleService) Find(filter map[string]any, data *model.Role) error {
	return common.DB.Where(filter).First(&data).Error
}

// Delete removes roles and associated Casbin policies
func (s RoleService) Delete(roleIds []uint) error {
	var roles []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	if err == nil {
		for _, role := range roles {
			rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("role deleted but failed to remove associated policies")
				}
			}
		}
	}
	return err
}

// GetRolesByIds gets roles by IDs
func (s RoleService) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// GetRoleMenusById gets role's menu permissions
func (s RoleService) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// UpdateRoleMenus updates role's menu permissions
func (s RoleService) UpdateRoleMenus(role *model.Role) error {
	return common.DB.Model(role).Association("Menus").Replace(role.Menus)
}

// UpdateRoleApis updates role's API permissions (delete all and recreate)
func (s RoleService) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("failed to load role policy")
	}
	rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("failed to update role API permissions")
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("failed to update role API permissions")
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("updated role API permissions but failed to reload policy")
	}
	return nil
}
