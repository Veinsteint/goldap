package request

// PendingUserListReq list pending users request
type PendingUserListReq struct {
	Username string `json:"username" form:"username"`
	Nickname string `json:"nickname" form:"nickname"`
	Mail     string `json:"mail" form:"mail"`
	Status   uint   `json:"status" form:"status"` // 0:pending, 1:approved, 2:rejected
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// PendingUserReviewReq review pending user request
type PendingUserReviewReq struct {
	ID            uint   `json:"id" validate:"required"`
	Status        uint   `json:"status" validate:"required,oneof=1 2"` // 1:approve, 2:reject
	ReviewRemark  string `json:"reviewRemark" validate:"max=255"`
	DepartmentId  []uint `json:"departmentId"`
	Departments   string `json:"departments"`
	RoleIds       []uint `json:"roleIds"`
	UIDNumber     uint   `json:"uidNumber"`
	GIDNumber     uint   `json:"gidNumber"`
	HomeDirectory string `json:"homeDirectory"`
	LoginShell    string `json:"loginShell"`
	Mobile        string `json:"mobile"`
	JobNumber     string `json:"jobNumber"`
	Position      string `json:"position"`
	PostalAddress string `json:"postalAddress"`
	AllowSudo     bool   `json:"allowSudo"`
	AllowSSHKey   bool   `json:"allowSSHKey"`
	SudoRules     string `json:"sudoRules"`
}

// PendingUserDeleteReq delete pending users request
type PendingUserDeleteReq struct {
	Ids []uint `json:"ids" validate:"required"`
}
