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

type UserPreConfigService struct{}

// List returns paginated user pre-config list
func (s UserPreConfigService) List(req *request.UserPreConfigListReq) ([]*model.UserPreConfig, error) {
	var list []*model.UserPreConfig
	db := common.DB.Model(&model.UserPreConfig{}).Order("created_at DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// ListAll returns all user pre-config records without pagination
func (s UserPreConfigService) ListAll() ([]*model.UserPreConfig, error) {
	var list []*model.UserPreConfig
	err := common.DB.Model(&model.UserPreConfig{}).Order("created_at DESC").Find(&list).Error
	return list, err
}

// Count returns total pre-config count
func (s UserPreConfigService) Count(req *request.UserPreConfigListReq) (int64, error) {
	var count int64
	db := common.DB.Model(&model.UserPreConfig{})

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}

	err := db.Count(&count).Error
	return count, err
}

// Add creates a new user pre-config
func (s UserPreConfigService) Add(data *model.UserPreConfig) error {
	return common.DB.Create(data).Error
}

// Update modifies a user pre-config
func (s UserPreConfigService) Update(data *model.UserPreConfig) error {
	// Use map to update all fields including zero values
	// Column names must match what's defined in the model's column tags
	return common.DB.Model(&model.UserPreConfig{}).Where("id = ?", data.ID).Updates(map[string]any{
		"username":       data.Username,
		"mail":           data.Mail,
		"uid_number":     data.UIDNumber,
		"gid_number":     data.GIDNumber,
		"department_id":  data.DepartmentId,
		"departments":    data.Departments,
		"mobile":         data.Mobile,
		"job_number":     data.JobNumber,
		"position":       data.Position,
		"postal_address": data.PostalAddress,
		"introduction":   data.Introduction,
		"home_directory": data.HomeDirectory,
		"login_shell":    data.LoginShell,
		"nickname":       data.Nickname,
		"given_name":     data.GivenName,
		"remark":         data.Remark,
		"is_used":        data.IsUsed,
	}).Error
}

// Delete removes user pre-configs by IDs
func (s UserPreConfigService) Delete(ids []uint) error {
	return common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.UserPreConfig{}).Error
}

// Find gets single pre-config by filter
func (s UserPreConfigService) Find(filter map[string]any, data *model.UserPreConfig) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist checks if pre-config exists
func (s UserPreConfigService) Exist(filter map[string]any) bool {
	var dataObj model.UserPreConfig
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// GetByUsername retrieves pre-config by username
func (s UserPreConfigService) GetByUsername(username string) (*model.UserPreConfig, error) {
	var config model.UserPreConfig
	err := common.DB.Where("username = ?", username).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &config, err
}

// MarkAsUsed marks a pre-config as used
func (s UserPreConfigService) MarkAsUsed(id uint) error {
	return common.DB.Model(&model.UserPreConfig{}).Where("id = ?", id).Update("is_used", true).Error
}

// SyncFromUser updates pre-config from user data (for consistency)
func (s UserPreConfigService) SyncFromUser(user *model.User) error {
	var config model.UserPreConfig
	err := common.DB.Where("username = ?", user.Username).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil // No pre-config exists, nothing to sync
	}
	if err != nil {
		return err
	}

	// Update pre-config with user's current values
	config.Mail = user.Mail
	config.UIDNumber = user.UIDNumber
	config.GIDNumber = user.GIDNumber
	config.DepartmentId = user.DepartmentId
	config.Departments = user.Departments
	config.Mobile = user.Mobile
	config.JobNumber = user.JobNumber
	config.Position = user.Position
	config.PostalAddress = user.PostalAddress
	config.Introduction = user.Introduction
	config.HomeDirectory = user.HomeDirectory
	config.LoginShell = user.LoginShell
	config.Nickname = user.Nickname
	config.GivenName = user.GivenName

	return common.DB.Model(&config).Updates(&config).Error
}

