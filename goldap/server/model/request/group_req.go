package request

// GroupListReq list groups request
type GroupListReq struct {
	GroupName string `json:"groupName" form:"groupName"`
	Remark    string `json:"remark" form:"remark"`
	PageNum   int    `json:"pageNum" form:"pageNum"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	SyncState uint   `json:"syncState" form:"syncState"`
}

// GroupListAllReq list all groups (no pagination)
type GroupListAllReq struct {
	GroupName          string `json:"groupName" form:"groupName"`
	GroupType          string `json:"groupType" form:"groupType"`
	Remark             string `json:"remark" form:"remark"`
	Source             string `json:"source" form:"source"`
	SourceDeptId       string `json:"sourceDeptId"`
	SourceDeptParentId string `json:"SourceDeptParentId"`
}

// GroupAddReq create group request
type GroupAddReq struct {
	GroupType string `json:"groupType" validate:"required,min=1,max=20"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	ParentId  uint   `json:"parentId" validate:"omitempty,min=0"`
	GIDNumber uint   `json:"gidNumber"` // POSIX GID number for posixGroup
	Remark    string `json:"remark" validate:"min=0,max=128"`
}

// GroupUpdateReq update group request
type GroupUpdateReq struct {
	ID        uint   `json:"id" form:"id" validate:"required"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	GIDNumber uint   `json:"gidNumber"` // POSIX GID number for posixGroup
	Remark    string `json:"remark" validate:"min=0,max=128"`
	IPRanges  string `json:"ipRanges"`
}

// GroupDeleteReq delete groups request
type GroupDeleteReq struct {
	GroupIds []uint `json:"groupIds" validate:"required"`
}

// GroupGetTreeReq get group tree request
type GroupGetTreeReq struct {
	GroupName string `json:"groupName" form:"groupName"`
	Remark    string `json:"remark" form:"remark"`
	PageNum   int    `json:"pageNum" form:"pageNum"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
}

// GroupAddUserReq add users to group
type GroupAddUserReq struct {
	GroupID uint   `json:"groupId" validate:"required"`
	UserIds []uint `json:"userIds" validate:"required"`
}

// GroupRemoveUserReq remove users from group
type GroupRemoveUserReq struct {
	GroupID uint   `json:"groupId" validate:"required"`
	UserIds []uint `json:"userIds" validate:"required"`
}

// UserInGroupReq query users in group
type UserInGroupReq struct {
	GroupID  uint   `json:"groupId" form:"groupId" validate:"required"`
	Nickname string `json:"nickname" form:"nickname"`
}

// UserNoInGroupReq query users not in group
type UserNoInGroupReq struct {
	GroupID  uint   `json:"groupId" form:"groupId" validate:"required"`
	Nickname string `json:"nickname" form:"nickname"`
}

// SyncOpenLdapDeptsReq sync LDAP departments
type SyncOpenLdapDeptsReq struct{}

// SyncSqlGrooupsReq sync MySQL groups to LDAP
type SyncSqlGrooupsReq struct {
	GroupIds []uint `json:"groupIds" validate:"required"`
}
