package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	// Initialize configuration
	cfg := config.NewServerConfig()

	agentRepo := db.NewAgentRepository(cfg.DB)
	userRepo := db.NewUserRepository(cfg.DB)
	commandQueue := queue.NewMemoryMultiQueue[models.Command]()

	gin.SetMode(gin.ReleaseMode)

	server := service.NewServer(agentRepo, commandQueue, cfg)

	if cfg.AgentHTTPConfig.Enabled {
		fmt.Println("enabled")
		agentAPIEngine := gin.Default()
		agentAPIHandler := api.NewAgentHandler(agentRepo, commandQueue, agentAPIEngine)
		httpTransport := transport.HTTPServerTransport(cfg.AgentHTTPConfig, agentAPIHandler)
		api.RegisterAgentRoutes(agentAPIHandler)
		server.AddTransport(httpTransport)
	}

	frontendEngine := gin.Default()
	frontendHandler := api.NewFrontendHandler(agentRepo, frontendEngine)
	api.RegisterFrontendRoutes(frontendHandler, cfg)
	api.RegisterAuthRoutes(frontendEngine, cfg, userRepo)
	api.StaticHandler(frontendEngine)
	frontendServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.FrontHTTPConfig.Port),
		Handler: frontendEngine,
	}

	// Start frontend server in a goroutine
	go func() {
		log.Printf("Starting frontend server on port %d", cfg.FrontHTTPConfig.Port)
		if err := frontendServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Frontend server error: %v", err)
		}
	}()

	// UDP transport (commented out, preserved for future use)
	// if cfg.UDPConfig.Enabled {
	// 	udpTransport := transport.NewUDPTransport(cfg.UDPConfig)
	// 	server.AddTransport(udpTransport)
	// }

	// Start the server for agent handling
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	log.Println("Server running. Press Ctrl+C to stop.")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	if err := frontendServer.Shutdown(context.Background()); err != nil {
		log.Printf("Frontend server shutdown error: %v", err)
	}
	server.Stop()
}
