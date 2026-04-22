package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitUserRoutes registers user management routes
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	user := r.Group("/user")
	user.Use(authMiddleware.MiddlewareFunc())
	user.Use(middleware.CasbinMiddleware())
	{
		user.GET("/info", controller.User.GetUserInfo)
		user.GET("/list", controller.User.List)
		user.POST("/add", controller.User.Add)
		user.POST("/update", controller.User.Update)
		user.POST("/delete", controller.User.Delete)
		user.POST("/changePwd", controller.User.ChangePwd)
		user.POST("/resetPassword", controller.User.ResetPassword)
		user.POST("/changeUserStatus", controller.User.ChangeUserStatus)

		user.POST("/syncOpenLdapUsers", controller.User.SyncOpenLdapUsers)
		user.POST("/syncSqlUsers", controller.User.SyncSqlUsers)

		// SSH key management
		user.GET("/ssh-keys", controller.SSHKey.GetSSHKeys)
		user.POST("/ssh-keys", controller.SSHKey.AddSSHKey)
		user.DELETE("/ssh-keys/:id", controller.SSHKey.DeleteSSHKey)

		// Pending user management
		user.GET("/pending/list", controller.PendingUser.List)
		user.POST("/pending/review", controller.PendingUser.Review)
		user.POST("/pending/delete", controller.PendingUser.Delete)
	}
	return r
}
