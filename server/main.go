package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/api"
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/service"
	"github.com/RabbITCybErSeC/BaconC2/server/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	logging.SetLevel(logging.LevelDebug)

	cfg := config.NewServerConfig()

	agentRepo := db.NewAgentRepository(cfg.DB)
	userRepo := db.NewUserRepository(cfg.DB)
	commandQueue := queue.NewMemoryMultiQueue[models.Command]()

	gin.SetMode(gin.ReleaseMode)

	server := service.NewServer(agentRepo, commandQueue, cfg)

	if cfg.AgentHTTPConfig.Enabled {
		agentAPIEngine := gin.Default()
		httpTransport := transport.NewHTTPServerTransport(agentRepo, commandQueue, cfg.AgentHTTPConfig, agentAPIEngine)
		server.AddTransport(httpTransport)
	}

	ginEngine := gin.Default()
	frontendHandler := api.NewFrontendHandler(agentRepo, ginEngine)
	generalHandler := api.NewGeneralApiHandler(agentRepo, ginEngine)

	api.RegisterFrontendRoutes(frontendHandler, cfg)
	api.RegisterApiRoutes(generalHandler, cfg)
	api.RegisterAuthRoutes(ginEngine, cfg, userRepo)
	api.StaticHandler(ginEngine)

	frontendServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.FrontHTTPConfig.Port),
		Handler: ginEngine,
	}

	go func() {
		logging.Info("Starting frontend server on port %d", cfg.FrontHTTPConfig.Port)
		if err := frontendServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Frontend server error: %v", err)
		}
	}()

	if err := server.Start(); err != nil {
		logging.Error("Failed to start services: %v", err)
		os.Exit(1)
	}

	logging.Info("Server running. Press Ctrl+C to stop.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logging.Info("Shutting down...")
	if err := frontendServer.Shutdown(context.Background()); err != nil {
		logging.Error("Frontend server shutdown error: %v", err)
	}
	server.Stop()
}
