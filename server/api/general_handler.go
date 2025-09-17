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
type GeneralApiHandler struct {
	agentRepository db.IAgentRepository
	commandQueue    queue.IServerCommandQueue
	engine          *gin.Engine
}

// NewAgentHandler initializes a new AgentHandler
func NewGeneralApiHandler(agentRepository db.IAgentRepository, commandQueue queue.IServerCommandQueue, engine *gin.Engine) *GeneralApiHandler {
	return &GeneralApiHandler{
		agentRepository: agentRepository,
		commandQueue:    commandQueue,
		engine:          engine,
	}
}

// GinEngine returns the Gin engine
func (h *GeneralApiHandler) GinEngine() *gin.Engine {
	return h.engine
}

// handleAddCommand handles adding a new command for an agent
func (h *GeneralApiHandler) handleAddCommand(c *gin.Context) {
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
			Status:  models.CommandStatusPending,
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
