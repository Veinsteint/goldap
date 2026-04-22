package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type MenuService struct{}

// Exist checks if menu exists
func (s MenuService) Exist(filter map[string]any) bool {
	var dataObj model.Menu
	err := common.DB.Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Count returns total menu count
func (s MenuService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Menu{}).Count(&count).Error
	return count, err
}

// Add creates a new menu
func (s MenuService) Add(menu *model.Menu) error {
	return common.DB.Create(menu).Error
}

// Update modifies menu
func (s MenuService) Update(menu *model.Menu) error {
	return common.DB.Model(&model.Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

// Find gets single menu by filter
func (s MenuService) Find(filter map[string]any, data *model.Menu) error {
	return common.DB.Where(filter).First(&data).Error
}

// List returns all menus ordered by sort
func (s MenuService) List() (menus []*model.Menu, err error) {
	err = common.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// ListUserMenus returns menus for given role IDs
func (s MenuService) ListUserMenus(roleIds []uint) (menus []*model.Menu, err error) {
	err = common.DB.Where("id IN (select menu_id as id from role_menus where role_id IN (?))", roleIds).Order("sort").Find(&menus).Error
	return menus, err
}

// Delete removes menus by IDs
func (s MenuService) Delete(menuIds []uint) error {
	return common.DB.Where("id IN (?)", menuIds).Select("Roles").Unscoped().Delete(&model.Menu{}).Error
}

// GetUserMenusByUserId returns accessible menus for user
func (s MenuService) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	var user model.User
	err := common.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range user.Roles {
		var userRole model.Role
		err := common.DB.Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		allRoleMenus = append(allRoleMenus, userRole.Menus...)
	}

	// Deduplicate menus
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// Filter active menus
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// GenMenuTree generates menu tree structure
func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)
	for _, m := range menus {
		if m.ParentId == parentId {
			m.Children = GenMenuTree(m.ID, menus)
			tree = append(tree, m)
		}
	}
	return tree
}
