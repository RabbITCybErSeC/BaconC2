package api

import (
	"net/http"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/server/store"
	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	agentStore store.IAgentStore
	engine     *gin.Engine
}

func NewFrontendHandler(agentStore store.IAgentStore, engine *gin.Engine) *FrontendHandler {
	return &FrontendHandler{
		agentStore: agentStore,
		engine:     engine,
	}
}

func (h *FrontendHandler) GinEngine() *gin.Engine {
	return h.engine
}

func (h *FrontendHandler) handleListAgents(c *gin.Context) {
	agentList, err := h.agentStore.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonAgents := make([]models.Agent, 0, len(agentList))
	for _, agent := range agentList {
		jsonAgents = append(jsonAgents, *&agent.BaseAgentModel)
	}

	c.JSON(http.StatusOK, jsonAgents)
}
