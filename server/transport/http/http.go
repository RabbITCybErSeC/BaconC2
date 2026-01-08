package http

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
	"github.com/RabbITCybErSeC/BaconC2/server/transport"
	"github.com/gin-gonic/gin"
)

const (
	ProtocolName          = "http"
	defaultBeaconInterval = 10
	shutdownTimeout       = 5 * time.Second
)

type HTTPServerTransport struct {
	agentRepository   db.IAgentRepository
	commandQueue      queue.IServerCommandQueue
	engine            *gin.Engine
	server            *http.Server
	httpConfig        config.AgentHTTPConfig
	pluginDataHandler *PluginDataHandler
}

func NewHTTPServerTransport(agentRepository db.IAgentRepository, commandQueue queue.IServerCommandQueue, httpConfig config.AgentHTTPConfig, engine *gin.Engine, pluginProvider transport.IPluginDataProvider) transport.ITransportProtocol {
	ht := &HTTPServerTransport{
		agentRepository:   agentRepository,
		commandQueue:      commandQueue,
		engine:            engine,
		httpConfig:        httpConfig,
		pluginDataHandler: NewPluginDataHandler(pluginProvider),
	}

	ht.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpConfig.Port),
		Handler: ht.engine,
	}
	ht.registerAgentRoutes()

	return ht
}

func (ht *HTTPServerTransport) registerAgentRoutes() {
	agentAPI := ht.engine.Group("/api/v1/agents")
	{
		agentAPI.Use(middleware.CorsMiddleware())
		agentAPI.POST("/register", ht.handleRegister)
		agentAPI.POST("/beacon", ht.handleBeacon)
		agentAPI.POST("/results", ht.handleCommandResult)

		agentAPI.GET("/plugins/metadata", ht.pluginDataHandler.HandlePluginMetadata)
		agentAPI.POST("/plugins/chunk", ht.pluginDataHandler.HandlePluginChunk)
	}
}

func (ht *HTTPServerTransport) GinEngine() *gin.Engine {
	return ht.engine
}

func (ht *HTTPServerTransport) requireAgentID(c *gin.Context) (string, bool) {
	agentID := c.Query("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return "", false
	}
	return agentID, true
}

func (ht *HTTPServerTransport) handleRegister(c *gin.Context) {
	var incomingAgent models.Agent
	if err := c.ShouldBindJSON(&incomingAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incomingAgent.Protocol = ProtocolName
	agent := local_models.ServerAgentModel{
		Agent:    incomingAgent,
		LastSeen: time.Now(),
		IsActive: true,
		Commands: []local_models.AgentCommand{},
	}

	if err := ht.agentRepository.SaveAgent(c.Request.Context(), &agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

func (ht *HTTPServerTransport) handleBeacon(c *gin.Context) {
	agentID, ok := ht.requireAgentID(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()

	if _, err := ht.agentRepository.GetAgent(ctx, agentID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := ht.agentRepository.UpdateLastSeen(ctx, agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update last seen: " + err.Error()})
		return
	}

	commands, err := ht.agentRepository.GetCommandsByStatus(ctx, agentID, models.CommandStatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(commands) > 0 {
		cmd := commands[0]

		if err := ht.agentRepository.UpdateCommandStatus(ctx, cmd.ID, models.CommandStatusSentToClient); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update command status"})
			return
		}

		c.JSON(http.StatusOK, models.HttpBeaconResponse{
			Command:        cmd.Command,
			Status:         models.CommandStatusSentToClient,
			NextBeacon:     defaultBeaconInterval,
			RequestResults: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.HttpBeaconResponse{
		Status:         models.CommandStatusAck,
		NextBeacon:     defaultBeaconInterval,
		RequestResults: false,
	})
}

func (ht *HTTPServerTransport) handleCommandResult(c *gin.Context) {
	agentID, ok := ht.requireAgentID(c)
	if !ok {
		return
	}

	var results []models.CommandResult
	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	if _, err := ht.agentRepository.GetAgent(ctx, agentID); err != nil {
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

		err := ht.agentRepository.UpdateCommandStatusWithResult(ctx, agentID, result.ID, result.Status, result.Output)
		if err != nil {
			log.Printf("Failed to save result for command %s from agent %s: %v", result.ID, agentID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process command result for " + result.ID})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (ht *HTTPServerTransport) Start() error {
	log.Printf("Starting HTTP transport on port %d", ht.httpConfig.Port)
	go func() {
		if err := ht.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	return nil
}

func (ht *HTTPServerTransport) Stop() error {
	if ht.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return ht.server.Shutdown(ctx)
}

func (ht *HTTPServerTransport) Name() string {
	return ProtocolName
}
