package request

// SystemConfigGetReq get system config request
type SystemConfigGetReq struct {
	Key string `json:"key" form:"key"`
}

// SystemConfigUpdateReq update system config request
type SystemConfigUpdateReq struct {
	PasswordSelfServiceEnabled bool   `json:"passwordSelfServiceEnabled"` // 密码自助服务开关
	RegistrationDisabled       bool   `json:"registrationDisabled"`       // 禁止注册开关
	SystemMaintenanceMode      bool   `json:"systemMaintenanceMode"`      // 系统维护中开关
	MaintenanceMessage         string `json:"maintenanceMessage"`         // 维护提示信息
}

