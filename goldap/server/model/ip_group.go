package model

import "gorm.io/gorm"

// IPGroup represents an IP-based access group
type IPGroup struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;unique;comment:'Group name'" json:"name"`
	Description string `gorm:"type:varchar(255);comment:'Description'" json:"description"`
	IPRanges    string `gorm:"type:text;not null;comment:'IP ranges (JSON array)'" json:"ipRanges"`
	Creator     string `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
}

// IPGroupUserPermission maps user permissions to IP groups
type IPGroupUserPermission struct {
	gorm.Model
	IPGroupID   uint    `gorm:"type:int;not null;index;comment:'IP group ID'" json:"ipGroupId"`
	UserID      uint    `gorm:"type:int;not null;index;comment:'User ID'" json:"userId"`
	AllowLogin  bool    `gorm:"type:tinyint(1);default:1;comment:'Allow login'" json:"allowLogin"`
	AllowSudo   bool    `gorm:"type:tinyint(1);default:0;comment:'Allow sudo'" json:"allowSudo"`
	AllowSSHKey bool    `gorm:"type:tinyint(1);default:1;comment:'Allow SSH key'" json:"allowSSHKey"`
	SudoRules   string  `gorm:"type:text;comment:'Sudo rules (JSON)'" json:"sudoRules"`
	IPGroup     IPGroup `gorm:"foreignKey:IPGroupID" json:"ipGroup,omitempty"`
	User        User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (g *IPGroup) SetName(name string) {
	g.Name = name
}

func (g *IPGroup) SetDescription(description string) {
	g.Description = description
}

func (g *IPGroup) SetIPRanges(ipRanges string) {
	g.IPRanges = ipRanges
}
