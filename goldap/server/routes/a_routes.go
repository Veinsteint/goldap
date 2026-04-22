package routes

import (
	"fmt"
	"net/http"
	"time"

	"goldap-server/config"
	_ "goldap-server/docs"
	"goldap-server/middleware"
	"goldap-server/public/common"
	"goldap-server/public/static"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRoutes initializes all application routes
func InitRoutes() *gin.Engine {
	// Set release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Use gin.New() instead of gin.Default() to disable default logging
	r := gin.New()
	r.Use(gin.Recovery()) // Only keep recovery middleware

	// Serve embedded static files
	r.Use(middleware.Serve("/", middleware.EmbedFolder(static.Static, "dist")))
	r.NoRoute(func(c *gin.Context) {
		data, err := static.Static.ReadFile("dist/index.html")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// Rate limiting middleware
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Operation log middleware
	r.Use(middleware.OperationLogMiddleware())

	// JWT authentication middleware
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		common.Log.Panicf("Failed to initialize JWT middleware: %v", err)
		panic(fmt.Sprintf("Failed to initialize JWT middleware: %v", err))
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes group
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)

	// Register routes
	InitBaseRoutes(apiGroup, authMiddleware)
	InitUserRoutes(apiGroup, authMiddleware)
	InitGroupRoutes(apiGroup, authMiddleware)
	InitRoleRoutes(apiGroup, authMiddleware)
	InitMenuRoutes(apiGroup, authMiddleware)
	InitApiRoutes(apiGroup, authMiddleware)
	InitOperationLogRoutes(apiGroup, authMiddleware)
	InitFieldRelationRoutes(apiGroup, authMiddleware)
	InitIPGroupRoutes(apiGroup, authMiddleware)
	InitIPGroupUserPermissionRoutes(apiGroup, authMiddleware)
	InitSudoRuleRoutes(apiGroup, authMiddleware)
	InitUserPreConfigRoutes(apiGroup, authMiddleware)
	InitSystemConfigRoutes(apiGroup, authMiddleware)
	return r
}
