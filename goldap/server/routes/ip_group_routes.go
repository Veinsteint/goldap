package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitIPGroupRoutes registers IP group routes
func InitIPGroupRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	ipGroup := r.Group("/ip-group")
	ipGroup.Use(authMiddleware.MiddlewareFunc())
	ipGroup.Use(middleware.CasbinMiddleware())
	{
		ipGroup.GET("", controller.IPGroup.List)
		ipGroup.POST("", controller.IPGroup.Add)
		ipGroup.PUT("", controller.IPGroup.Update)
		ipGroup.DELETE("", controller.IPGroup.Delete)
	}
	return r
}

// InitIPGroupUserPermissionRoutes registers IP group user permission routes
func InitIPGroupUserPermissionRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	permission := r.Group("/ip-group-user-permission")
	permission.Use(authMiddleware.MiddlewareFunc())
	permission.Use(middleware.CasbinMiddleware())
	{
		permission.GET("", controller.IPGroupUserPermission.List)
		permission.POST("", controller.IPGroupUserPermission.Add)
		permission.PUT("", controller.IPGroupUserPermission.Update)
		permission.DELETE("", controller.IPGroupUserPermission.Delete)
	}
	return r
}

// InitSudoRuleRoutes registers sudo rule routes
func InitSudoRuleRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	sudoRule := r.Group("/sudo-rule")
	sudoRule.Use(authMiddleware.MiddlewareFunc())
	sudoRule.Use(middleware.CasbinMiddleware())
	{
		sudoRule.GET("", controller.SudoRule.List)
		sudoRule.POST("", controller.SudoRule.Add)
		sudoRule.PUT("", controller.SudoRule.Update)
		sudoRule.DELETE("", controller.SudoRule.Delete)
	}
	return r
}
