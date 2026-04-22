package model

import "gorm.io/gorm"

// UserPreConfig stores pre-configured user information
// When approving new users, these values are used as defaults if the username matches
type UserPreConfig struct {
	gorm.Model
	Username      string `gorm:"column:username;type:varchar(50);not null;unique;comment:'Username (matching key)'" json:"username"`
	Mail          string `gorm:"column:mail;type:varchar(100);comment:'Pre-configured email'" json:"mail"`
	UIDNumber     uint   `gorm:"column:uid_number;type:int;comment:'Pre-configured UID'" json:"uidNumber"`
	GIDNumber     uint   `gorm:"column:gid_number;type:int;comment:'Pre-configured GID'" json:"gidNumber"`
	DepartmentId  string `gorm:"column:department_id;type:varchar(100);comment:'Pre-configured department IDs'" json:"departmentId"`
	Departments   string `gorm:"column:departments;type:varchar(512);comment:'Pre-configured department names'" json:"departments"`
	Mobile        string `gorm:"column:mobile;type:varchar(15);comment:'Pre-configured mobile'" json:"mobile"`
	JobNumber     string `gorm:"column:job_number;type:varchar(20);comment:'Pre-configured job number'" json:"jobNumber"`
	Position      string `gorm:"column:position;type:varchar(128);comment:'Pre-configured position'" json:"position"`
	PostalAddress string `gorm:"column:postal_address;type:varchar(255);comment:'Pre-configured address'" json:"postalAddress"`
	Introduction  string `gorm:"column:introduction;type:varchar(255);comment:'Pre-configured introduction'" json:"introduction"`
	HomeDirectory string `gorm:"column:home_directory;type:varchar(255);comment:'Pre-configured home directory'" json:"homeDirectory"`
	LoginShell    string `gorm:"column:login_shell;type:varchar(255);comment:'Pre-configured login shell'" json:"loginShell"`
	Nickname      string `gorm:"column:nickname;type:varchar(50);comment:'Pre-configured nickname'" json:"nickname"`
	GivenName     string `gorm:"column:given_name;type:varchar(50);comment:'Pre-configured given name'" json:"givenName"`
	Creator       string `gorm:"column:creator;type:varchar(20);comment:'Creator'" json:"creator"`
	Remark        string `gorm:"column:remark;type:varchar(255);comment:'Remark'" json:"remark"`
	IsUsed        bool   `gorm:"column:is_used;type:tinyint(1);default:0;comment:'Whether this config has been used'" json:"isUsed"`
}

