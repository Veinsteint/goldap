package request

// GroupUserPermissionAddReq add group user permission
type GroupUserPermissionAddReq struct {
	GroupID     uint   `json:"groupId" validate:"required"`
	UserID      uint   `json:"userId" validate:"required"`
	AllowSudo   bool   `json:"allowSudo"`
	AllowSSHKey bool   `json:"allowSSHKey"`
	SudoRules   string `json:"sudoRules"`
}

// GroupUserPermissionUpdateReq update group user permission
type GroupUserPermissionUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	AllowSudo   bool   `json:"allowSudo"`
	AllowSSHKey bool   `json:"allowSSHKey"`
	SudoRules   string `json:"sudoRules"`
}

// GroupUserPermissionDeleteReq delete group user permission
type GroupUserPermissionDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}

// GroupUserPermissionListReq list group user permissions
type GroupUserPermissionListReq struct {
	GroupID uint `json:"groupId"`
	UserID  uint `json:"userId"`
}
