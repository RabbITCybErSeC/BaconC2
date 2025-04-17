package api

// RegisterAgentRoutes sets up routes for AgentHandler
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

// RegisterFrontendRoutes sets up routes for FrontendHandler
func RegisterFrontendRoutes(frontendHandler *FrontendHandler) {
	frontendAPI := frontendHandler.engine.Group("/api/frontend")
	{
		frontendAPI.Use(CorsMiddleware())
		frontendAPI.GET("/agents", frontendHandler.handleListAgents)
	}
}
