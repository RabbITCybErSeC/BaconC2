package plugins

import (
	"fmt"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

type PluginTransferManager struct {
	compiler *PluginCompiler
}

func NewPluginTransferManager(compiler *PluginCompiler) *PluginTransferManager {
	return &PluginTransferManager{
		compiler: compiler,
	}
}

func (m *PluginTransferManager) CreatePluginLoadCommand(pluginName string) (*models.Command, error) {
	pluginPath := m.compiler.GetPluginPath(pluginName)
	hash, err := plugins.CalculateHash(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	cmd := models.NewCommand("plugin_install", models.CommandTypeInternal, pluginName, hash)

	logging.Info("Created plugin install command for '%s' (hash: %s)", pluginName, hash)
	return cmd, nil
}

func (m *PluginTransferManager) CreatePluginUnloadCommand(pluginName string) *models.Command {
	return models.NewCommand("plugin_unload", models.CommandTypeInternal, pluginName)
}

func (m *PluginTransferManager) CreatePluginListCommand() *models.Command {
	return models.NewCommand("plugin_list", models.CommandTypeInternal)
}

func (m *PluginTransferManager) CreatePluginStatusCommand(pluginName string) *models.Command {
	return models.NewCommand("plugin_status", models.CommandTypeInternal, pluginName)
}

func (m *PluginTransferManager) CreatePluginInfoCommand(pluginName string) *models.Command {
	return models.NewCommand("plugin_info", models.CommandTypeInternal, pluginName)
}
