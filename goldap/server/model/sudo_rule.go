package model

import "gorm.io/gorm"

// SudoRule represents a sudo permission rule
type SudoRule struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;unique;comment:'Rule name'" json:"name"`
	Description string `gorm:"type:varchar(255);comment:'Description'" json:"description"`
	User        string `gorm:"type:varchar(255);comment:'User or group (e.g. %group1 or user1)'" json:"user"`
	Host        string `gorm:"type:varchar(255);default:ALL;comment:'Host (e.g. ALL or 192.168.1.0/24)'" json:"host"`
	Command     string `gorm:"type:text;not null;comment:'Allowed commands'" json:"command"`
	RunAsUser   string `gorm:"type:varchar(255);default:ALL;comment:'Run as user'" json:"runAsUser"`
	RunAsGroup  string `gorm:"type:varchar(255);default:ALL;comment:'Run as group'" json:"runAsGroup"`
	Options     string `gorm:"type:varchar(255);comment:'Options (e.g. NOPASSWD, NOEXEC)'" json:"options"`
	Creator     string `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
}

func (s *SudoRule) SetName(name string) {
	s.Name = name
}

func (s *SudoRule) SetDescription(description string) {
	s.Description = description
}
