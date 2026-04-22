package model

import (
	"gorm.io/gorm"
)

// Menu represents a navigation menu item
type Menu struct {
	gorm.Model
	Name       string  `gorm:"type:varchar(50);comment:'Menu name (i18n key)'" json:"name"`
	Title      string  `gorm:"type:varchar(50);comment:'Menu title'" json:"title"`
	Icon       string  `gorm:"type:varchar(50);comment:'Menu icon'" json:"icon"`
	Path       string  `gorm:"type:varchar(100);comment:'Route path'" json:"path"`
	Redirect   string  `gorm:"type:varchar(100);comment:'Redirect path'" json:"redirect"`
	Component  string  `gorm:"type:varchar(100);comment:'Component path'" json:"component"`
	Sort       uint    `gorm:"type:int(3);default:999;comment:'Sort order (1-999)'" json:"sort"`
	Status     uint    `gorm:"type:tinyint(1);default:1;comment:'Status: 1=Active, 2=Disabled'" json:"status"`
	Hidden     uint    `gorm:"type:tinyint(1);default:2;comment:'Hidden: 1=Yes, 2=No'" json:"hidden"`
	NoCache    uint    `gorm:"type:tinyint(1);default:2;comment:'No cache: 1=Yes, 2=No'" json:"noCache"`
	AlwaysShow uint    `gorm:"type:tinyint(1);default:2;comment:'Always show: 1=Yes, 2=No'" json:"alwaysShow"`
	Breadcrumb uint    `gorm:"type:tinyint(1);default:1;comment:'Breadcrumb: 1=Show, 2=Hide'" json:"breadcrumb"`
	ActiveMenu string  `gorm:"type:varchar(100);comment:'Active menu path'" json:"activeMenu"`
	ParentId   uint    `gorm:"default:0;comment:'Parent ID (0=root)'" json:"parentId"`
	Creator    string  `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	Children   []*Menu `gorm:"-" json:"children"`
	Roles      []*Role `gorm:"many2many:role_menus;" json:"roles"`
}
