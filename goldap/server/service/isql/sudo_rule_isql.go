package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type SudoRuleService struct{}

// Add creates sudo rule
func (s SudoRuleService) Add(rule *model.SudoRule) error {
	return common.DB.Create(rule).Error
}

// Update modifies sudo rule
func (s SudoRuleService) Update(rule *model.SudoRule) error {
	return common.DB.Model(rule).Updates(rule).Error
}

// Delete removes sudo rule
func (s SudoRuleService) Delete(id uint) error {
	return common.DB.Delete(&model.SudoRule{}, id).Error
}

// Find gets sudo rule by filter
func (s SudoRuleService) Find(filter map[string]any, data *model.SudoRule) error {
	return common.DB.Where(filter).First(data).Error
}

// List returns all sudo rules
func (s SudoRuleService) List() ([]*model.SudoRule, error) {
	var list []*model.SudoRule
	err := common.DB.Order("created_at DESC").Find(&list).Error
	return list, err
}

// GetByID gets sudo rule by ID
func (s SudoRuleService) GetByID(id uint) (*model.SudoRule, error) {
	var rule model.SudoRule
	err := common.DB.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Exist checks if sudo rule exists
func (s SudoRuleService) Exist(filter map[string]any) bool {
	var dataObj model.SudoRule
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
