package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitGroupRoutes registers group management routes
func InitGroupRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	group := r.Group("/group")
	group.Use(authMiddleware.MiddlewareFunc())
	group.Use(middleware.CasbinMiddleware())
	{
		group.GET("/list", controller.Group.List)
		group.GET("/tree", controller.Group.GetTree)
		group.POST("/add", controller.Group.Add)
		group.POST("/update", controller.Group.Update)
		group.POST("/delete", controller.Group.Delete)
		group.POST("/adduser", controller.Group.AddUser)
		group.POST("/removeuser", controller.Group.RemoveUser)
		group.GET("/useringroup", controller.Group.UserInGroup)
		group.GET("/usernoingroup", controller.Group.UserNoInGroup)
		group.POST("/syncOpenLdapDepts", controller.Group.SyncOpenLdapDepts)
		group.POST("/syncSqlGroups", controller.Group.SyncSqlGroups)
	}

	// Group user permission routes
	permission := r.Group("/group-user-permission")
	permission.Use(authMiddleware.MiddlewareFunc())
	permission.Use(middleware.CasbinMiddleware())
	{
		permission.GET("", controller.GroupUserPermission.List)
		permission.POST("", controller.GroupUserPermission.Add)
		permission.PUT("", controller.GroupUserPermission.Update)
		permission.DELETE("", controller.GroupUserPermission.Delete)
	}

	return r
}
