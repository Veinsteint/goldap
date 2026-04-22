package model

import "gorm.io/gorm"

// GroupUserPermission maps user permissions to groups
type GroupUserPermission struct {
	gorm.Model
	GroupID     uint   `gorm:"type:bigint unsigned;not null;index;comment:'Group ID'" json:"groupId"`
	UserID      uint   `gorm:"type:bigint unsigned;not null;index;comment:'User ID'" json:"userId"`
	AllowSudo   bool   `gorm:"type:tinyint(1);default:0;comment:'Allow sudo'" json:"allowSudo"`
	AllowSSHKey bool   `gorm:"type:tinyint(1);default:1;comment:'Allow SSH key upload'" json:"allowSSHKey"`
	SudoRules   string `gorm:"type:text;comment:'Sudo rules (JSON)'" json:"sudoRules"`
	Group       Group  `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
