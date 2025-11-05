package http

import (
	"net/http"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/server/transport"
	"github.com/gin-gonic/gin"
)

type PluginDataHandler struct {
	provider transport.IPluginDataProvider
}

func NewPluginDataHandler(provider transport.IPluginDataProvider) *PluginDataHandler {
	return &PluginDataHandler{
		provider: provider,
	}
}

func (h *PluginDataHandler) HandlePluginMetadata(c *gin.Context) {
	pluginName := c.Query("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name required"})
		return
	}

	metadata, err := h.provider.GetPluginMetadata(pluginName)
	if err != nil {
		logging.Error("Failed to get plugin metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

func (h *PluginDataHandler) HandlePluginChunk(c *gin.Context) {
	var req models.PluginChunkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.PluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name required"})
		return
	}

	if req.ChunkIndex < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunk index"})
		return
	}

	response, err := h.provider.GetPluginChunk(req.PluginName, req.ChunkIndex)
	if err != nil {
		logging.Error("Failed to get plugin chunk: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
