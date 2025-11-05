package transport

import (
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type IPluginFetcher interface {
	FetchPluginMetadata(pluginName string) (*models.PluginTransferMetadata, error)
	FetchPluginChunk(pluginName string, chunkIndex int) (*models.PluginChunkResponse, error)
}
