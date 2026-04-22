package logic

import (
	"fmt"
	"strconv"

	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type SystemConfigLogic struct{}

// Get retrieves system configuration
func (l SystemConfigLogic) Get(c *gin.Context, req any) (data any, rspError any) {
	_ = c
	_ = req

	// Get config values
	passwordSelfServiceEnabled := isql.SystemConfig.GetValue("password_self_service_enabled", "true") == "true"
	registrationDisabled := isql.SystemConfig.GetValue("registration_disabled", "false") == "true"
	systemMaintenanceMode := isql.SystemConfig.GetValue("system_maintenance_mode", "false") == "true"
	maintenanceMessage := isql.SystemConfig.GetValue("maintenance_message", "系统正在升级维护中，请稍后再试")

	return response.SystemConfigRsp{
		PasswordSelfServiceEnabled: passwordSelfServiceEnabled,
		RegistrationDisabled:       registrationDisabled,
		SystemMaintenanceMode:      systemMaintenanceMode,
		MaintenanceMessage:         maintenanceMessage,
	}, nil
}

// Update updates system configuration
func (l SystemConfigLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SystemConfigUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Update config values
	if err := isql.SystemConfig.SetValue("password_self_service_enabled", strconv.FormatBool(r.PasswordSelfServiceEnabled), "密码自助服务开关"); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update password self service config: %v", err))
	}

	if err := isql.SystemConfig.SetValue("registration_disabled", strconv.FormatBool(r.RegistrationDisabled), "禁止注册开关"); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update registration disabled config: %v", err))
	}

	if err := isql.SystemConfig.SetValue("system_maintenance_mode", strconv.FormatBool(r.SystemMaintenanceMode), "系统维护中开关"); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update system maintenance mode config: %v", err))
	}

	if err := isql.SystemConfig.SetValue("maintenance_message", r.MaintenanceMessage, "维护提示信息"); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update maintenance message config: %v", err))
	}

	return map[string]string{
		"message": "System configuration updated successfully",
	}, nil
}

