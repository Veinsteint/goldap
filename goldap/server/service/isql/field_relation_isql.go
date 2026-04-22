package isql

import (
	"errors"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type FieldRelationService struct{}

// Exist checks if field relation exists
func (s FieldRelationService) Exist(filter map[string]any) bool {
	var dataObj model.FieldRelation
	err := common.DB.Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Count returns total count
func (s FieldRelationService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.FieldRelation{}).Count(&count).Error
	return count, err
}

// Add creates field relation
func (s FieldRelationService) Add(fieldRelation *model.FieldRelation) error {
	return common.DB.Create(fieldRelation).Error
}

// Update modifies field relation
func (s FieldRelationService) Update(fieldRelation *model.FieldRelation) error {
	return common.DB.Model(&model.FieldRelation{}).Where("id = ?", fieldRelation.ID).Updates(fieldRelation).Error
}

// Find gets single field relation
func (s FieldRelationService) Find(filter map[string]any, data *model.FieldRelation) error {
	return common.DB.Where(filter).First(&data).Error
}

// List returns all field relations
func (s FieldRelationService) List() (fieldRelations []*model.FieldRelation, err error) {
	err = common.DB.Find(&fieldRelations).Error
	return fieldRelations, err
}

// Delete removes field relations by IDs
func (s FieldRelationService) Delete(fieldRelationIds []uint) error {
	return common.DB.Where("id IN (?)", fieldRelationIds).Unscoped().Delete(&model.FieldRelation{}).Error
}
