package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"goldap-server/controller"
)

// LoginHandler placeholder for Swagger documentation
// @Summary Login
// @Description User login (add Bearer + token for password encryption)
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param data body request.RegisterAndLoginReq true "Login credentials"
// @Success 200 {object} response.ResponseBody
// @Router /base/login [post]
func LoginHandler() {}

// LogoutHandler placeholder for Swagger documentation
// @Summary Logout
// @Description User logout
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/logout [post]
func LogoutHandler() {}

// RefreshHandler placeholder for Swagger documentation
// @Summary Refresh Token
// @Description Refresh JWT token
// @Tags Base
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.ResponseBody
// @Router /base/refreshToken [post]
func RefreshHandler() {}

// InitBaseRoutes registers base routes (no authentication required)
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", controller.Demo)
		base.GET("encryptpwd", controller.Base.EncryptPasswd)
		base.GET("decryptpwd", controller.Base.DecryptPasswd)
		base.POST("/login", authMiddleware.LoginHandler)
		base.POST("/logout", authMiddleware.LogoutHandler)
		base.POST("/refreshToken", authMiddleware.RefreshHandler)
		base.POST("/sendcode", controller.Base.SendCode)
		base.POST("/changePwd", controller.Base.ChangePwd)
		base.POST("/register", controller.PendingUser.Register)
		base.GET("/dashboard", controller.Base.Dashboard)
		base.GET("/validUsernames", controller.UserPreConfig.GetValidUsernames)
		base.GET("/registrationMode", controller.UserPreConfig.GetRegistrationMode)
		base.GET("/systemConfig", controller.Base.GetSystemConfig)
	}
	return r
}
