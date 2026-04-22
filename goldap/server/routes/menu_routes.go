package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitMenuRoutes registers menu management routes
func InitMenuRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	menu := r.Group("/menu")
	menu.Use(authMiddleware.MiddlewareFunc())
	menu.Use(middleware.CasbinMiddleware())
	{
		menu.GET("/tree", controller.Menu.GetTree)
		menu.GET("/access/tree", controller.Menu.GetAccessTree)
		menu.POST("/add", controller.Menu.Add)
		menu.POST("/update", controller.Menu.Update)
		menu.POST("/delete", controller.Menu.Delete)
	}
	return r
}
