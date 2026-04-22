package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type IPGroupUserPermissionService struct{}

// Add creates IP group user permission
func (s IPGroupUserPermissionService) Add(permission *model.IPGroupUserPermission) error {
	return common.DB.Create(permission).Error
}

// Update modifies IP group user permission
func (s IPGroupUserPermissionService) Update(permission *model.IPGroupUserPermission) error {
	return common.DB.Model(permission).Updates(permission).Error
}

// Delete removes IP group user permission
func (s IPGroupUserPermissionService) Delete(id uint) error {
	return common.DB.Delete(&model.IPGroupUserPermission{}, id).Error
}

// Find gets IP group user permission by filter
func (s IPGroupUserPermissionService) Find(filter map[string]any, data *model.IPGroupUserPermission) error {
	return common.DB.Where(filter).Preload("IPGroup").Preload("User").First(data).Error
}

// GetByUserID gets all permissions for user
func (s IPGroupUserPermissionService) GetByUserID(userID uint) ([]*model.IPGroupUserPermission, error) {
	var list []*model.IPGroupUserPermission
	err := common.DB.Where("user_id = ?", userID).Preload("IPGroup").Preload("User").Find(&list).Error
	return list, err
}

// GetByIPGroupID gets all permissions for IP group
func (s IPGroupUserPermissionService) GetByIPGroupID(ipGroupID uint) ([]*model.IPGroupUserPermission, error) {
	var list []*model.IPGroupUserPermission
	err := common.DB.Where("ip_group_id = ?", ipGroupID).Preload("IPGroup").Preload("User").Find(&list).Error
	return list, err
}

// GetByUserAndIP gets permission by user ID and client IP
func (s IPGroupUserPermissionService) GetByUserAndIP(userID uint, clientIP string) (*model.IPGroupUserPermission, error) {
	var permissions []*model.IPGroupUserPermission
	err := common.DB.Where("user_id = ? AND allow_login = ?", userID, true).
		Preload("IPGroup").Preload("User").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	
	if len(permissions) > 0 {
		return permissions[0], nil
	}
	return nil, gorm.ErrRecordNotFound
}

// Exist checks if permission exists
func (s IPGroupUserPermissionService) Exist(filter map[string]any) bool {
	var dataObj model.IPGroupUserPermission
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
