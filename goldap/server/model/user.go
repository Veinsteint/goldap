package model

import "gorm.io/gorm"

// User represents a system user
type User struct {
	gorm.Model
	Username      string  `gorm:"type:varchar(50);not null;unique;comment:'Username'" json:"username"`
	Password      string  `gorm:"type:text;not null;comment:'Password'" json:"password"`
	Nickname      string  `gorm:"type:varchar(50);comment:'Display name'" json:"nickname"`
	GivenName     string  `gorm:"type:varchar(50);comment:'Given name'" json:"givenName"`
	Mail          string  `gorm:"type:varchar(100);not null;comment:'Email'" json:"mail"`
	JobNumber     string  `gorm:"type:varchar(20);comment:'Employee number'" json:"jobNumber"`
	Mobile        string  `gorm:"type:varchar(15);comment:'Mobile phone'" json:"mobile"`
	Avatar        string  `gorm:"type:varchar(255);comment:'Avatar URL'" json:"avatar"`
	PostalAddress string  `gorm:"type:varchar(255);comment:'Address'" json:"postalAddress"`
	Departments   string  `gorm:"type:varchar(512);comment:'Departments'" json:"departments"`
	Position      string  `gorm:"type:varchar(128);comment:'Position'" json:"position"`
	Introduction  string  `gorm:"type:varchar(255);comment:'Introduction'" json:"introduction"`
	Status        uint    `gorm:"type:tinyint(1);default:1;comment:'Status: 1=Active, 2=Inactive'" json:"status"`
	Creator       string  `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	Source        string  `gorm:"type:varchar(50);comment:'Source: ldap, platform'" json:"source"`
	DepartmentId  string  `gorm:"type:varchar(100);comment:'Department IDs'" json:"departmentId"`
	Roles         []*Role `gorm:"many2many:user_roles" json:"roles"`
	SourceUserId  string  `gorm:"type:varchar(100);not null;comment:'Third-party user ID'" json:"sourceUserId"`
	SourceUnionId string  `gorm:"type:varchar(100);not null;comment:'Third-party union ID'" json:"sourceUnionId"`
	UserDN        string  `gorm:"type:varchar(255);not null;comment:'LDAP DN'" json:"userDn"`
	SyncState     uint    `gorm:"type:tinyint(1);default:1;comment:'Sync state: 1=Synced, 2=Pending'" json:"syncState"`

	// Unix user fields
	UIDNumber     uint   `gorm:"column:uid_number;type:int;unique;comment:'Unix UID'" json:"uidNumber"`
	GIDNumber     uint   `gorm:"column:gid_number;type:int;default:0;comment:'Unix GID'" json:"gidNumber"`
	HomeDirectory string `gorm:"column:home_directory;type:varchar(255);comment:'Home directory'" json:"homeDirectory"`
	LoginShell    string `gorm:"column:login_shell;type:varchar(255);default:/bin/bash;comment:'Login shell'" json:"loginShell"`
	Gecos         string `gorm:"column:gecos;type:varchar(255);comment:'GECOS field'" json:"gecos"`
}

func (u *User) SetUserName(userName string) {
	u.Username = userName
}

func (u *User) SetNickName(nickName string) {
	u.Nickname = nickName
}

func (u *User) SetGivenName(givenName string) {
	u.GivenName = givenName
}

func (u *User) SetMail(mail string) {
	u.Mail = mail
}

func (u *User) SetJobNumber(jobNum string) {
	u.JobNumber = jobNum
}

func (u *User) SetMobile(mobile string) {
	u.Mobile = mobile
}

func (u *User) SetAvatar(avatar string) {
	u.Avatar = avatar
}

func (u *User) SetPostalAddress(address string) {
	u.PostalAddress = address
}

func (u *User) SetPosition(position string) {
	u.Position = position
}

func (u *User) SetIntroduction(desc string) {
	u.Introduction = desc
}

func (u *User) SetSourceUserId(sourceUserId string) {
	u.SourceUserId = sourceUserId
}

func (u *User) SetSourceUnionId(sourceUnionId string) {
	u.SourceUnionId = sourceUnionId
}
