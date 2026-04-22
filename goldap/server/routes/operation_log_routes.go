package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitOperationLogRoutes registers operation log routes
func InitOperationLogRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	operationLog := r.Group("/log")
	operationLog.Use(authMiddleware.MiddlewareFunc())
	operationLog.Use(middleware.CasbinMiddleware())
	{
		operationLog.GET("/operation/list", controller.OperationLog.List)
		operationLog.POST("/operation/delete", controller.OperationLog.Delete)
		operationLog.DELETE("/operation/clean", controller.OperationLog.Clean)
	}
	return r
}
