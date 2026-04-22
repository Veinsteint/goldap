package model

import "gorm.io/gorm"

// Role represents a user role for RBAC
type Role struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(20);not null;unique" json:"name"`
	Keyword string  `gorm:"type:varchar(20);not null;unique" json:"keyword"`
	Remark  string  `gorm:"type:varchar(100);comment:'Description'" json:"remark"`
	Status  uint    `gorm:"type:tinyint(1);default:1;comment:'Status: 1=Active, 2=Disabled'" json:"status"`
	Sort    uint    `gorm:"type:int(3);default:999;comment:'Sort order (1=Super Admin)'" json:"sort"`
	Creator string  `gorm:"type:varchar(20);" json:"creator"`
	Users   []*User `gorm:"many2many:user_roles" json:"users"`
	Menus   []*Menu `gorm:"many2many:role_menus;" json:"menus"`
}
