package routes

import (
	"goldap-server/controller"
	"goldap-server/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitFieldRelationRoutes registers field relation routes
func InitFieldRelationRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	fieldRelation := r.Group("/fieldrelation")
	fieldRelation.Use(authMiddleware.MiddlewareFunc())
	fieldRelation.Use(middleware.CasbinMiddleware())
	{
		fieldRelation.POST("/add", controller.FieldRelation.Add)
		fieldRelation.GET("/list", controller.FieldRelation.List)
		fieldRelation.POST("/update", controller.FieldRelation.Update)
		fieldRelation.POST("/delete", controller.FieldRelation.Delete)
	}
	return r
}
