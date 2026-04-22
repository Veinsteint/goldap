package request

// UserAddReq create user request
type UserAddReq struct {
	Username        string `json:"username" validate:"required,min=2,max=50"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Nickname        string `json:"nickname" validate:"required,min=0,max=50"`
	GivenName       string `json:"givenName" validate:"min=0,max=50"`
	Mail            string `json:"mail" validate:"required,email,max=100"`
	JobNumber       string `json:"jobNumber" validate:"min=0,max=20"`
	PostalAddress   string `json:"postalAddress" validate:"min=0,max=255"`
	Departments     string `json:"departments" validate:"min=0,max=512"`
	Position        string `json:"position" validate:"min=0,max=128"`
	Mobile          string `json:"mobile" validate:"checkMobile"`
	Avatar          string `json:"avatar"`
	Introduction    string `json:"introduction" validate:"min=0,max=255"`
	Status          uint   `json:"status" validate:"oneof=1 2"`
	DepartmentId    []uint `json:"departmentId"`
	Source          string `json:"source" validate:"min=0,max=50"`
	RoleIds         []uint `json:"roleIds"`
	UIDNumber       uint   `json:"uidNumber"`
	GIDNumber       uint   `json:"gidNumber"`
	HomeDirectory   string `json:"homeDirectory"`
	LoginShell      string `json:"loginShell"`
	Gecos           string `json:"gecos"`
}

// DingUserAddReq third-party user create request
type DingUserAddReq struct {
	Username        string `json:"username" validate:"required,min=2,max=50"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Nickname        string `json:"nickname" validate:"required,min=0,max=50"`
	GivenName       string `json:"givenName" validate:"min=0,max=50"`
	Mail            string `json:"mail" validate:"required,min=0,max=100"`
	JobNumber       string `json:"jobNumber" validate:"min=0,max=20"`
	PostalAddress   string `json:"postalAddress" validate:"min=0,max=255"`
	Departments     string `json:"departments" validate:"min=0,max=512"`
	Position        string `json:"position" validate:"min=0,max=128"`
	Mobile          string `json:"mobile" validate:"checkMobile"`
	Avatar          string `json:"avatar"`
	Introduction    string `json:"introduction" validate:"min=0,max=255"`
	Status          uint   `json:"status" validate:"oneof=1 2"`
	DepartmentId    []uint `json:"departmentId" validate:"required"`
	Source          string `json:"source" validate:"min=0,max=50"`
	RoleIds         []uint `json:"roleIds" validate:"required"`
	SourceUserId    string `json:"sourceUserId"`
	SourceUnionId   string `json:"sourceUnionId"`
}

// UserUpdateReq update user request
type UserUpdateReq struct {
	ID            uint   `json:"id" validate:"required"`
	Username      string `json:"username" validate:"required,min=2,max=50"`
	Nickname      string `json:"nickname" validate:"min=0,max=20"`
	GivenName     string `json:"givenName" validate:"min=0,max=50"`
	Mail          string `json:"mail" validate:"min=0,max=100"`
	JobNumber     string `json:"jobNumber" validate:"min=0,max=20"`
	PostalAddress string `json:"postalAddress" validate:"min=0,max=255"`
	Departments   string `json:"departments" validate:"min=0,max=512"`
	Position      string `json:"position" validate:"min=0,max=128"`
	Mobile        string `json:"mobile" validate:"checkMobile"`
	Avatar        string `json:"avatar"`
	Introduction  string `json:"introduction" validate:"min=0,max=255"`
	DepartmentId  []uint `json:"departmentId" validate:"required"`
	Source        string `json:"source" validate:"min=0,max=50"`
	RoleIds       []uint `json:"roleIds" validate:"required"`
	UIDNumber     uint   `json:"uidNumber"`
	GIDNumber     uint   `json:"gidNumber"`
	HomeDirectory string `json:"homeDirectory"`
	LoginShell    string `json:"loginShell"`
	Gecos         string `json:"gecos"`
}

// UserDeleteReq batch delete users request
type UserDeleteReq struct {
	UserIds []uint `json:"userIds" validate:"required"`
}

// UserChangePwdReq change password request
type UserChangePwdReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

// UserResetPasswordReq reset password request
type UserResetPasswordReq struct {
	Username string `json:"username" validate:"required"`
}

// UserChangeUserStatusReq change user status request
type UserChangeUserStatusReq struct {
	ID     uint `json:"id" validate:"required"`
	Status uint `json:"status" validate:"oneof=1 2"`
}

// UserGetUserInfoReq get user info request
type UserGetUserInfoReq struct{}

// SyncOpenLdapUserReq sync LDAP users request
type SyncOpenLdapUserReq struct{}

// SyncSqlUserReq sync MySQL users to LDAP
type SyncSqlUserReq struct {
	UserIds []uint `json:"userIds"`
}

// UserListReq list users request
type UserListReq struct {
	Username     string `json:"username" form:"username"`
	Mobile       string `json:"mobile" form:"mobile"`
	Nickname     string `json:"nickname" form:"nickname"`
	GivenName    string `json:"givenName" form:"givenName"`
	DepartmentId []uint `json:"departmentId" form:"departmentId"`
	Status       uint   `json:"status" form:"status"`
	SyncState    uint   `json:"syncState" form:"syncState"`
	PageNum      int    `json:"pageNum" form:"pageNum"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
}

// RegisterAndLoginReq login request
type RegisterAndLoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// UserRegisterReq user registration request (requires admin approval)
type UserRegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RealName string `json:"realName" validate:"required,max=50"`
	Remark   string `json:"remark" validate:"max=200"`
}
