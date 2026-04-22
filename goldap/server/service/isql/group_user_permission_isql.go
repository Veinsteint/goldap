package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type GroupUserPermissionService struct{}

// Add creates group user permission
func (s GroupUserPermissionService) Add(permission *model.GroupUserPermission) error {
	return common.DB.Create(permission).Error
}

// Update modifies group user permission
func (s GroupUserPermissionService) Update(permission *model.GroupUserPermission) error {
	return common.DB.Model(permission).Updates(permission).Error
}

// Delete removes group user permission
func (s GroupUserPermissionService) Delete(id uint) error {
	return common.DB.Delete(&model.GroupUserPermission{}, id).Error
}

// Find gets group user permission by filter
func (s GroupUserPermissionService) Find(filter map[string]any, data *model.GroupUserPermission) error {
	return common.DB.Where(filter).Preload("Group").Preload("User").First(data).Error
}

// GetByUserID gets all permissions for user
func (s GroupUserPermissionService) GetByUserID(userID uint) ([]*model.GroupUserPermission, error) {
	var list []*model.GroupUserPermission
	err := common.DB.Where("user_id = ?", userID).Preload("Group").Preload("User").Find(&list).Error
	return list, err
}

// GetByGroupID gets all permissions for group
func (s GroupUserPermissionService) GetByGroupID(groupID uint) ([]*model.GroupUserPermission, error) {
	var list []*model.GroupUserPermission
	err := common.DB.Where("group_id = ?", groupID).Preload("Group").Preload("User").Find(&list).Error
	return list, err
}

// GetByGroupAndUser gets permission by group and user
func (s GroupUserPermissionService) GetByGroupAndUser(groupID, userID uint) (*model.GroupUserPermission, error) {
	var permission model.GroupUserPermission
	err := common.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		Preload("Group").Preload("User").First(&permission).Error
	return &permission, err
}

// Exist checks if permission exists
func (s GroupUserPermissionService) Exist(filter map[string]any) bool {
	var dataObj model.GroupUserPermission
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
