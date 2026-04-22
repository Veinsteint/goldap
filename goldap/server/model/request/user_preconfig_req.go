package request

// UserPreConfigListReq request for listing user pre-configs
type UserPreConfigListReq struct {
	Username string `json:"username" form:"username"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// UserPreConfigAddReq request for adding user pre-config
type UserPreConfigAddReq struct {
	Username      string `json:"username" form:"username" validate:"required,min=2,max=50"`
	Mail          string `json:"mail" form:"mail"`
	UIDNumber     uint   `json:"uidNumber" form:"uidNumber"`
	GIDNumber     uint   `json:"gidNumber" form:"gidNumber"`
	DepartmentId  string `json:"departmentId" form:"departmentId"`
	Departments   string `json:"departments" form:"departments"`
	Mobile        string `json:"mobile" form:"mobile"`
	JobNumber     string `json:"jobNumber" form:"jobNumber"`
	Position      string `json:"position" form:"position"`
	PostalAddress string `json:"postalAddress" form:"postalAddress"`
	Introduction  string `json:"introduction" form:"introduction"`
	HomeDirectory string `json:"homeDirectory" form:"homeDirectory"`
	LoginShell    string `json:"loginShell" form:"loginShell"`
	Nickname      string `json:"nickname" form:"nickname"`
	GivenName     string `json:"givenName" form:"givenName"`
	Remark        string `json:"remark" form:"remark"`
}

// UserPreConfigUpdateReq request for updating user pre-config
type UserPreConfigUpdateReq struct {
	ID            uint   `json:"id" form:"id" validate:"required"`
	Username      string `json:"username" form:"username" validate:"required,min=2,max=50"`
	Mail          string `json:"mail" form:"mail"`
	UIDNumber     uint   `json:"uidNumber" form:"uidNumber"`
	GIDNumber     uint   `json:"gidNumber" form:"gidNumber"`
	DepartmentId  string `json:"departmentId" form:"departmentId"`
	Departments   string `json:"departments" form:"departments"`
	Mobile        string `json:"mobile" form:"mobile"`
	JobNumber     string `json:"jobNumber" form:"jobNumber"`
	Position      string `json:"position" form:"position"`
	PostalAddress string `json:"postalAddress" form:"postalAddress"`
	Introduction  string `json:"introduction" form:"introduction"`
	HomeDirectory string `json:"homeDirectory" form:"homeDirectory"`
	LoginShell    string `json:"loginShell" form:"loginShell"`
	Nickname      string `json:"nickname" form:"nickname"`
	GivenName     string `json:"givenName" form:"givenName"`
	Remark        string `json:"remark" form:"remark"`
}

// UserPreConfigDeleteReq request for deleting user pre-configs
type UserPreConfigDeleteReq struct {
	Ids []uint `json:"ids" form:"ids" validate:"required"`
}

// UserPreConfigGetByUsernameReq request for getting pre-config by username
type UserPreConfigGetByUsernameReq struct {
	Username string `json:"username" form:"username" validate:"required"`
}

