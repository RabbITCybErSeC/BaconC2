package api

import (
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
	agentList, err := h.agentRepository.GetAll()
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

func (h *FrontendHandler) handleGetAgentByID(c *gin.Context) {
	id := c.Param("id")
	agent, err := h.agentRepository.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}
