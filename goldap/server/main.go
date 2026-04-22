package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goldap-server/config"
	"goldap-server/logic"
	"goldap-server/middleware"
	"goldap-server/public/common"
	"goldap-server/routes"
	"goldap-server/service/isql"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// @title GoLDAP Server API
// @version 1.0
// @description OpenLDAP management system API
// @host 127.0.0.1:8150
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.InitConfig()
	common.InitLogger()

	common.Log.Infof("GoLDAP Server %s (build: %s, commit: %s)", Version, BuildTime, GitCommit)

	common.InitDB()
	common.InitLDAP()

	if err := common.EnsureSudoSchema(); err != nil {
		common.Log.Warnf("sudo schema check failed: %v", err)
	}

	if err := logic.ClearLDAP(); err != nil {
		common.Log.Warnf("LDAP clear failed: %v", err)
	}

	common.InitCasbinEnforcer()
	common.InitValidate()
	common.InitData()

	// Async sync MySQL users to LDAP
	go syncDataOnStartup()

	// Start operation log workers
	for i := 0; i < 3; i++ {
		go isql.OperationLog.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Conf.System.Port),
		Handler:      routes.InitRoutes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		common.Log.Infof("Server starting on port %d", config.Conf.System.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("Server failed: %v", err)
		}
	}()

	logic.InitCron()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	common.Log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatalf("Server shutdown failed: %v", err)
	}
	common.Log.Info("Server stopped")
}

func syncDataOnStartup() {
	time.Sleep(2 * time.Second)

	if err := logic.SearchUserDiff(); err != nil {
		common.Log.Errorf("User sync failed: %v", err)
	}

	if err := logic.SyncGroupMembersFromDB(); err != nil {
		common.Log.Errorf("Group sync failed: %v", err)
	}
}
