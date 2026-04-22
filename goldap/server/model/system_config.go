package model

import "gorm.io/gorm"

// SystemConfig represents system configuration settings
type SystemConfig struct {
	gorm.Model
	Key         string `gorm:"column:key;type:varchar(100);not null;unique;comment:'Config key'" json:"key"`
	Value       string `gorm:"type:varchar(255);comment:'Config value'" json:"value"`
	Description string `gorm:"type:varchar(255);comment:'Config description'" json:"description"`
}

// TableName returns the table name for SystemConfig
func (SystemConfig) TableName() string {
	return "system_configs"
}

