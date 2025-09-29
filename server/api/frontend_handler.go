package api

import (
	"errors"
	"net/http"

	"github.com/RabbITCybErSeC/BaconC2/server/db"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	agentRepository db.IAgentRepository
	engine          *gin.Engine
}

func NewFrontendHandler(agentRepository db.IAgentRepository, engine *gin.Engine) *FrontendHandler {
	return &FrontendHandler{
		agentRepository: agentRepository,
		engine:          engine,
	}
}

func (h *FrontendHandler) GinEngine() *gin.Engine {
	return h.engine
}

func (h *FrontendHandler) handleListAgents(c *gin.Context) {
	agentList, err := h.agentRepository.GetAllAgents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonAgents := make([]local_models.ServerAgentModel, 0, len(agentList))
	for _, agent := range agentList {
		jsonAgents = append(jsonAgents, agent)
	}

	c.JSON(http.StatusOK, jsonAgents)
}

func (h *GeneralApiHandler) handleGetCommandResult(c *gin.Context) {
	commandID := c.Param("commandId")
	if commandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command ID is required"})
		return
	}

	result, err := h.agentRepository.GetCommandResult(c.Request.Context(), commandID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Command result not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve command result: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FrontendHandler) handleGetAgentByID(c *gin.Context) {
	id := c.Param("id")
	agent, err := h.agentRepository.GetAgent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}
