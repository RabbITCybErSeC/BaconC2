package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	agentRepository db.IAgentRepository
	engine          *gin.Engine
	pluginDir       string
}

func NewFrontendHandler(agentRepository db.IAgentRepository, engine *gin.Engine, pluginDir string) *FrontendHandler {
	return &FrontendHandler{
		agentRepository: agentRepository,
		engine:          engine,
		pluginDir:       pluginDir,
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

func (h *FrontendHandler) handleGetPlugins(c *gin.Context) {
	plugins, err := h.getLoadedPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plugins: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, plugins)
}

// getLoadedPlugins scans the plugin directory and returns metadata for all .so files
func (h *FrontendHandler) getLoadedPlugins() ([]plugins.PluginInfo, error) {
	if h.pluginDir == "" {
		return []plugins.PluginInfo{}, nil
	}

	// Check if directory exists
	if _, err := os.Stat(h.pluginDir); os.IsNotExist(err) {
		return []plugins.PluginInfo{}, nil
	}

	var pluginInfos []plugins.PluginInfo

	// Walk through the plugin directory
	err := filepath.Walk(h.pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.so files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".so") {
			return nil
		}

		pluginName := strings.TrimSuffix(info.Name(), ".so")
		jsonPath := filepath.Join(h.pluginDir, pluginName+".json")

		// Try to read metadata from JSON file first
		if jsonData, err := os.ReadFile(jsonPath); err == nil {
			var pluginInfo plugins.PluginInfo
			if err := json.Unmarshal(jsonData, &pluginInfo); err == nil {
				// Update file size and hash from actual file
				pluginInfo.FileSize = info.Size()
				hash, _ := calculateFileHash(path)
				pluginInfo.Hash = hash
				pluginInfo.Metadata.Hash = hash
				pluginInfo.Metadata.LoadedAt = info.ModTime()
				pluginInfo.FilePath = path
				pluginInfo.Metadata.FilePath = path
				
				pluginInfos = append(pluginInfos, pluginInfo)
				return nil
			}
		}

		// Fallback: Create generic metadata if JSON doesn't exist or is invalid
		hash, err := calculateFileHash(path)
		if err != nil {
			return fmt.Errorf("failed to calculate hash for %s: %w", path, err)
		}

		metadata := plugins.PluginMetadata{
			Name:        pluginName,
			Version:     "unknown",
			Author:      "unknown",
			Description: "Plugin loaded from filesystem",
			LoadedAt:    info.ModTime(),
			Status:      plugins.PluginStatusLoaded,
			FilePath:    path,
			Hash:        hash,
		}

		pluginInfo := plugins.PluginInfo{
			Metadata: metadata,
			FilePath: path,
			FileSize: info.Size(),
			Hash:     hash,
		}

		pluginInfos = append(pluginInfos, pluginInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan plugin directory: %w", err)
	}

	return pluginInfos, nil
}

func calculateFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
