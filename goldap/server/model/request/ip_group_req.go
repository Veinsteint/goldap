package request

// IPGroupAddReq create IP group request
type IPGroupAddReq struct {
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"max=255"`
	IPRanges    []string `json:"ipRanges" validate:"required,min=1"`
}

// IPGroupUpdateReq update IP group request
type IPGroupUpdateReq struct {
	ID          uint     `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"max=255"`
	IPRanges    []string `json:"ipRanges" validate:"required,min=1"`
}

// IPGroupDeleteReq delete IP group request
type IPGroupDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}

// IPGroupUserPermissionAddReq add IP group user permission
type IPGroupUserPermissionAddReq struct {
	IPGroupID   uint   `json:"ipGroupId" validate:"required"`
	UserID      uint   `json:"userId" validate:"required"`
	AllowLogin  bool   `json:"allowLogin"`
	AllowSudo   bool   `json:"allowSudo"`
	AllowSSHKey bool   `json:"allowSSHKey"`
	SudoRules   string `json:"sudoRules"`
}

// IPGroupUserPermissionUpdateReq update IP group user permission
type IPGroupUserPermissionUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	AllowLogin  bool   `json:"allowLogin"`
	AllowSudo   bool   `json:"allowSudo"`
	AllowSSHKey bool   `json:"allowSSHKey"`
	SudoRules   string `json:"sudoRules"`
}

// IPGroupUserPermissionDeleteReq delete IP group user permission
type IPGroupUserPermissionDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}

// IPGroupUserPermissionListReq list IP group user permissions
type IPGroupUserPermissionListReq struct {
	IPGroupID uint `json:"ipGroupId"`
	UserID    uint `json:"userId"`
}
