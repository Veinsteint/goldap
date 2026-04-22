package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type IPGroupService struct{}

// Add creates IP group
func (s IPGroupService) Add(ipGroup *model.IPGroup) error {
	return common.DB.Create(ipGroup).Error
}

// Update modifies IP group
func (s IPGroupService) Update(ipGroup *model.IPGroup) error {
	return common.DB.Model(ipGroup).Updates(ipGroup).Error
}

// Delete removes IP group and associated permissions
func (s IPGroupService) Delete(id uint) error {
	common.DB.Where("ip_group_id = ?", id).Delete(&model.IPGroupUserPermission{})
	return common.DB.Delete(&model.IPGroup{}, id).Error
}

// Find gets IP group by filter
func (s IPGroupService) Find(filter map[string]any, data *model.IPGroup) error {
	return common.DB.Where(filter).First(data).Error
}

// List returns all IP groups
func (s IPGroupService) List() ([]*model.IPGroup, error) {
	var list []*model.IPGroup
	err := common.DB.Order("created_at DESC").Find(&list).Error
	return list, err
}

// Exist checks if IP group exists
func (s IPGroupService) Exist(filter map[string]any) bool {
	var dataObj model.IPGroup
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// GetByID gets IP group by ID
func (s IPGroupService) GetByID(id uint) (*model.IPGroup, error) {
	var ipGroup model.IPGroup
	err := common.DB.First(&ipGroup, id).Error
	if err != nil {
		return nil, err
	}
	return &ipGroup, nil
}
