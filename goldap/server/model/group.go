package model

import (
	"gorm.io/gorm"
)

// Group represents a user group or organizational unit
type Group struct {
	gorm.Model
	GroupName          string   `gorm:"type:varchar(128);comment:'Group name'" json:"groupName"`
	Remark             string   `gorm:"type:varchar(128);comment:'Description'" json:"remark"`
	Creator            string   `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	GroupType          string   `gorm:"type:varchar(20);comment:'Type: cn, ou, posix'" json:"groupType"`
	GIDNumber          uint     `gorm:"column:gid_number;type:int;default:0;comment:'POSIX GID number'" json:"gidNumber"`
	Users              []*User  `gorm:"many2many:group_users" json:"users"`
	ParentId           uint     `gorm:"default:0;comment:'Parent ID (0=root)'" json:"parentId"`
	SourceDeptId       string   `gorm:"type:varchar(100);comment:'Source department ID'" json:"sourceDeptId"`
	Source             string   `gorm:"type:varchar(20);comment:'Source: ldap, platform'" json:"source"`
	SourceDeptParentId string   `gorm:"type:varchar(100);comment:'Source parent department ID'" json:"sourceDeptParentId"`
	SourceUserNum      int      `gorm:"default:0;comment:'User count from source'" json:"source_user_num"`
	Children           []*Group `gorm:"-" json:"children"`
	GroupDN            string   `gorm:"type:varchar(255);not null;comment:'LDAP DN'" json:"groupDn"`
	SyncState          uint     `gorm:"type:tinyint(1);default:1;comment:'Sync state: 1=Synced, 2=Pending'" json:"syncState"`
	IPRanges           string   `gorm:"type:text;comment:'IP ranges (JSON array)'" json:"ipRanges"`
}

func (g *Group) SetGroupName(groupName string) {
	g.GroupName = groupName
}

func (g *Group) SetRemark(remark string) {
	g.Remark = remark
}

func (g *Group) SetSourceDeptId(sourceDeptId string) {
	g.SourceDeptId = sourceDeptId
}

func (g *Group) SetSourceDeptParentId(sourceDeptParentId string) {
	g.SourceDeptParentId = sourceDeptParentId
}
