package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// SendCode sends verification code to email
// @Summary Send Verification Code
// @Description Send verification code to specified email
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseSendCodeReq true "Send code request"
// @Success 200 {object} response.ResponseBody
// @Router /base/sendcode [post]
func (m *BaseController) SendCode(c *gin.Context) {
	req := new(request.BaseSendCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.SendCode(c, req)
	})
}

// ChangePwd changes password via email verification
// @Summary Change Password via Email
// @Description Change password using email verification code
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseChangePwdReq true "Change password request"
// @Success 200 {object} response.ResponseBody
// @Router /base/changePwd [post]
func (m *BaseController) ChangePwd(c *gin.Context) {
	req := new(request.BaseChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ChangePwd(c, req)
	})
}

// Dashboard returns system dashboard data
// @Summary Get Dashboard Data
// @Description Get system dashboard overview data
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/dashboard [get]
func (m *BaseController) Dashboard(c *gin.Context) {
	req := new(request.BaseDashboardReq)
	Run(c, req, func() (any, any) {
		return logic.Base.Dashboard(c, req)
	})
}

// EncryptPasswd encrypts plaintext password
// @Summary Encrypt Password
// @Description Encrypt plaintext password
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "Plaintext password to encrypt"
// @Success 200 {object} response.ResponseBody
// @Router /base/encryptpwd [get]
func (m *BaseController) EncryptPasswd(c *gin.Context) {
	req := new(request.EncryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.EncryptPasswd(c, req)
	})
}

// DecryptPasswd decrypts encrypted password
// @Summary Decrypt Password
// @Description Decrypt encrypted password to plaintext
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "Encrypted password to decrypt"
// @Success 200 {object} response.ResponseBody
// @Router /base/decryptpwd [get]
func (m *BaseController) DecryptPasswd(c *gin.Context) {
	req := new(request.DecryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.DecryptPasswd(c, req)
	})
}

// GetSystemConfig retrieves system configuration (public endpoint)
// @Summary Get System Configuration
// @Description Get system configuration settings (public endpoint)
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/systemConfig [get]
func (m *BaseController) GetSystemConfig(c *gin.Context) {
	req := new(request.SystemConfigGetReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetSystemConfig(c, req)
	})
}
