package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitUserPreConfigRoutes registers user pre-config routes
func InitUserPreConfigRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	preconfig := r.Group("/user/preconfig")
	preconfig.Use(authMiddleware.MiddlewareFunc())
	preconfig.Use(middleware.CasbinMiddleware())
	{
		preconfig.GET("/list", controller.UserPreConfig.List)
		preconfig.POST("/add", controller.UserPreConfig.Add)
		preconfig.POST("/update", controller.UserPreConfig.Update)
		preconfig.POST("/delete", controller.UserPreConfig.Delete)
		preconfig.GET("/getByUsername", controller.UserPreConfig.GetByUsername)
		preconfig.POST("/syncUsers", controller.UserPreConfig.SyncExistingUsers)
	}
	return r
}

