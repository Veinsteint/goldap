package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type BaseLogic struct{}

// SendCode sends verification code to email
func (l BaseLogic) SendCode(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseSendCodeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	user := new(model.User)
	if err := isql.User.Find(tools.H{"mail": r.Mail}, user); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("user not found: %s", err.Error()))
	}

	if user.Status != 1 || user.SyncState != 1 {
		return nil, tools.NewMySqlError(fmt.Errorf("user is inactive or not synced"))
	}

	if err := tools.SendCode([]string{r.Mail}); err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to send email: %s", err.Error()))
	}

	return nil, nil
}

// ChangePwd resets user password via email verification
func (l BaseLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if !isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorError(fmt.Errorf("email not found"))
	}

	cacheCode, ok := tools.VerificationCodeCache.Get(r.Mail)
	if !ok {
		return nil, tools.NewValidatorError(fmt.Errorf("verification code expired"))
	}

	if cacheCode != r.Code {
		return nil, tools.NewValidatorError(fmt.Errorf("invalid verification code"))
	}

	user := new(model.User)
	if err := isql.User.Find(tools.H{"mail": r.Mail}, user); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("user not found: %s", err.Error()))
	}

	newpass, err := ildap.User.NewPwd(user.Username, user.UserDN)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to generate password: %s", err.Error()))
	}

	if err := tools.SendMail([]string{user.Mail}, newpass); err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("failed to send email: %s", err.Error()))
	}

	if err := isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newpass)); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update password: %s", err.Error()))
	}

	return nil, nil
}

// Dashboard returns dashboard statistics
func (l BaseLogic) Dashboard(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseDashboardReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	userCount, _ := isql.User.Count()
	groupCount, _ := isql.Group.Count()
	roleCount, _ := isql.Role.Count()
	menuCount, _ := isql.Menu.Count()
	apiCount, _ := isql.Api.Count()
	logCount, _ := isql.OperationLog.Count()

	rst := []*response.DashboardList{
		{DataType: "user", DataName: "Users", DataCount: userCount, Icon: "people", Path: "#/personnel/user"},
		{DataType: "group", DataName: "Groups", DataCount: groupCount, Icon: "peoples", Path: "#/personnel/group"},
		{DataType: "role", DataName: "Roles", DataCount: roleCount, Icon: "eye-open", Path: "#/system/role"},
		{DataType: "menu", DataName: "Menus", DataCount: menuCount, Icon: "tree-table", Path: "#/system/menu"},
		{DataType: "api", DataName: "APIs", DataCount: apiCount, Icon: "tree", Path: "#/system/api"},
		{DataType: "log", DataName: "Logs", DataCount: logCount, Icon: "documentation", Path: "#/log/operation-log"},
	}

	return rst, nil
}

// EncryptPasswd encrypts password
func (l BaseLogic) EncryptPasswd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.EncryptPasswdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return tools.NewGenPasswd(r.Passwd), nil
}

// DecryptPasswd decrypts password
func (l BaseLogic) DecryptPasswd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.DecryptPasswdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return tools.NewParPasswd(r.Passwd), nil
}

// GetSystemConfig retrieves system configuration (public endpoint)
func (l BaseLogic) GetSystemConfig(c *gin.Context, req any) (data any, rspError any) {
	_ = c
	_ = req

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
