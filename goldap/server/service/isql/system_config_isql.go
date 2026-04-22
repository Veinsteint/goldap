package isql

import (
	"errors"
	"fmt"

	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type SystemConfigService struct{}

// GetByKey retrieves a system config by key
func (s SystemConfigService) GetByKey(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := common.DB.Where("`key` = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// GetValue retrieves a system config value by key, returns defaultValue if not found
func (s SystemConfigService) GetValue(key string, defaultValue string) string {
	config, err := s.GetByKey(key)
	if err != nil || config == nil {
		return defaultValue
	}
	return config.Value
}

// SetValue sets a system config value by key, creates if not exists
func (s SystemConfigService) SetValue(key string, value string, description string) error {
	var config model.SystemConfig
	err := common.DB.Where("`key` = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new config
			config = model.SystemConfig{
				Key:         key,
				Value:       value,
				Description: description,
			}
			return common.DB.Create(&config).Error
		}
		return err
	}
	// Update existing config
	config.Value = value
	if description != "" {
		config.Description = description
	}
	return common.DB.Save(&config).Error
}

// GetAll retrieves all system configs
func (s SystemConfigService) GetAll() ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	err := common.DB.Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get system configs: %v", err)
	}
	return configs, nil
}

