package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitSSHKeyRoutes registers SSH key routes (under /user group)
func InitSSHKeyRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	user := r.Group("/user")
	user.Use(authMiddleware.MiddlewareFunc())
	user.Use(middleware.CasbinMiddleware())
	{
		user.GET("/ssh-keys", controller.SSHKey.GetSSHKeys)
		user.POST("/ssh-keys", controller.SSHKey.AddSSHKey)
		user.DELETE("/ssh-keys/:id", controller.SSHKey.DeleteSSHKey)
	}
	return r
}
