package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitSystemConfigRoutes registers system configuration routes
func InitSystemConfigRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	system := r.Group("/system")
	system.Use(authMiddleware.MiddlewareFunc())
	system.Use(middleware.CasbinMiddleware())
	{
		config := system.Group("/config")
		{
			config.GET("/get", controller.SystemConfig.Get)
			config.POST("/update", controller.SystemConfig.Update)
		}
	}
	return r
}

