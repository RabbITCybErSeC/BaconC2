package api

import (
	"errors"
	"net/http"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent-related operations
type GeneralApiHandler struct {
	agentRepository db.IAgentRepository
	engine          *gin.Engine
}

// NewAgentHandler initializes a new AgentHandler
func NewGeneralApiHandler(agentRepository db.IAgentRepository, engine *gin.Engine) *GeneralApiHandler {
	return &GeneralApiHandler{
		agentRepository: agentRepository,
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

	if rawCmd.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'type' is required and cannot be empty"})
		return
	}

	agentID := c.Param("agentId")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	_, err := h.agentRepository.GetAgent(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	cmd := models.NewCommand(rawCmd.Command, rawCmd.Type, rawCmd.Args...)
	agentCmdPtr := local_models.NewAgentCommand(agentID, *cmd)

	if err := h.agentRepository.SaveCommand(c.Request.Context(), agentCmdPtr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued", "id": agentCmdPtr.ID})
}

func (h *GeneralApiHandler) handleGetAllAgentCommands(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID is required"})
		return
	}

	_, err := h.agentRepository.GetAgent(c.Request.Context(), agentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify agent: " + err.Error()})
		return
	}

	// A limit of 0 means get all commands (based on the existing repository implementation).
	commands, err := h.agentRepository.GetCommands(c.Request.Context(), agentID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve agent commands: " + err.Error()})
		return
	}

	if commands == nil {
		c.JSON(http.StatusOK, []local_models.AgentCommand{})
		return
	}

	c.JSON(http.StatusOK, commands)
}

func (h *GeneralApiHandler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Backend is running",
	})
}
