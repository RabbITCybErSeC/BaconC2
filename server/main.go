package main

import (
	"log"

	"github.com/RabbITCybErSeC/BaconC2/server/api"
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/service"
	"github.com/RabbITCybErSeC/BaconC2/server/store"
	"github.com/RabbITCybErSeC/BaconC2/server/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewServerConfig()

	agentRepo := db.NewAgentRepository(cfg.DB)
	agentStore := store.NewAgentStore(agentRepo)
	commandQueue := queue.NewMemoryCommandQueue()

	gin.SetMode(gin.ReleaseMode) // Set to release mode for production
	apiHandler := api.NewHandler(agentStore, commandQueue)

	server := service.NewServer(agentStore, commandQueue, cfg)

	if cfg.HTTPConfig.Enabled {
		httpTransport := transport.NewHTTPTransport(cfg.HTTPConfig, apiHandler)
		server.AddTransport(httpTransport)
	}
	// if cfg.UDPConfig.Enabled {
	// 	udpTransport := transport.NewUDPTransport(cfg.UDPConfig)
	// 	server.AddTransport(udpTransport)
	// }

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server running. Press Ctrl+C to stop.")
	select {} // Block forever
}
