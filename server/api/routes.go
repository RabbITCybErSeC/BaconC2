package api

import (
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/gin-gonic/gin"
)

func RegisterAgentRoutes(agentHandler *AgentHandler) {
	agentAPI := agentHandler.engine.Group("/api/agents")
	{
		agentAPI.Use(CorsMiddleware())
		agentAPI.POST("/register", agentHandler.handleRegister)
		agentAPI.POST("/beacon", agentHandler.handleBeacon)
		agentAPI.POST("/result", agentHandler.handleCommandResult)
		agentAPI.POST("/command", agentHandler.handleAddCommand)
	}
}

func RegisterFrontendRoutes(frontendHandler *FrontendHandler, config *config.ServerConfig) {
	frontendAPI := frontendHandler.engine.Group("/api/frontend")
	{
		frontendAPI.Use(CorsMiddleware())
		frontendAPI.Use(JWTMiddleware(config))
		frontendAPI.GET("/agents", frontendHandler.handleListAgents)
	}
}

func RegisterAuthRoutes(engine *gin.Engine, config *config.ServerConfig, userRepo db.UserRepositoryInterface) {
	authHandler := NewAuthHandler(config, userRepo)
	authAPI := engine.Group("/api/auth")
	{
		authAPI.Use(CorsMiddleware())
		authAPI.POST("/login", authHandler.handleLogin)
	}
}
