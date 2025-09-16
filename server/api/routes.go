package api

import (
	"fmt" // Replace with your actual module path
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/RabbITCybErSeC/BaconC2/pkg/middleware"
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"

	"github.com/RabbITCybErSeC/BaconC2/becongui"
	"github.com/gin-gonic/gin"
)

func RegisterAgentRoutes(agentHandler *AgentHandler) {
	agentAPI := agentHandler.engine.Group("/api/agents")
	{
		agentAPI.Use(middleware.CorsMiddleware())
		agentAPI.POST("/register", agentHandler.handleRegister)
		agentAPI.POST("/beacon", agentHandler.handleBeacon)
		agentAPI.POST("/result", agentHandler.handleCommandResult)
		agentAPI.POST("/command", agentHandler.handleAddCommand)
	}
}

func RegisterFrontendRoutes(frontendHandler *FrontendHandler, config *config.ServerConfig) {
	frontendAPI := frontendHandler.engine.Group("/api/frontend")
	{
		frontendAPI.Use(middleware.CorsMiddleware())
		frontendAPI.Use(JWTMiddleware(config))
		frontendAPI.GET("/agents", frontendHandler.handleListAgents)
	}
}

func RegisterAuthRoutes(engine *gin.Engine, config *config.ServerConfig, userRepo db.UserRepositoryInterface) {
	authHandler := NewAuthHandler(config, userRepo)
	authAPI := engine.Group("/api/auth")
	{
		authAPI.Use(middleware.CorsMiddleware())
		authAPI.POST("/login", authHandler.handleLogin)
	}
}

func StaticHandler(engine *gin.Engine) {
	// Create a sub-filesystem for the embedded dist folder
	dist, err := fs.Sub(becongui.Dist, "dist")
	if err != nil {
		panic(err) // Handle error appropriately in production
	}
	fileServer := http.FileServer(http.FS(dist))

	engine.Use(func(c *gin.Context) {
		// Skip API routes
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			// Check if the requested file exists
			path := strings.TrimPrefix(c.Request.URL.Path, "/")
			_, err := fs.Stat(dist, path)
			if os.IsNotExist(err) {
				// If the file does not exist, serve index.html for SPA routing
				fmt.Println("File not found, serving index.html")
				c.Request.URL.Path = "/"
			} else {
				// Serve other static files
				fmt.Println("Serving static file:", c.Request.URL.Path)
			}

			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}
