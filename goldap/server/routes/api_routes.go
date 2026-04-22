package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitApiRoutes registers API management routes
func InitApiRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	api.Use(middleware.CasbinMiddleware())
	{
		api.GET("/tree", controller.Api.GetTree)
		api.GET("/list", controller.Api.List)
		api.POST("/add", controller.Api.Add)
		api.POST("/update", controller.Api.Update)
		api.POST("/delete", controller.Api.Delete)
	}
	return r
}
