package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

// PluginService handles plugin lifecycle and integration with the command system
// It provides the high-level API for plugin operations
type PluginService struct {
	store           IPluginStore
	commandRegistry *command_handler.CommandHandlerRegistry
	agentState      command_handler.IAgentState
	mu              sync.RWMutex
	config          map[string]interface{}
}

// NewPluginService creates a new plugin service with Lua-based engine
func NewPluginService(
	pluginDir string,
	commandRegistry *command_handler.CommandHandlerRegistry,
	agentState command_handler.IAgentState,
) *PluginService {
	engine := NewLuaEngine(pluginDir)
	store := NewPluginStore(engine)

	return &PluginService{
		store:           store,
		commandRegistry: commandRegistry,
		agentState:      agentState,
		config:          make(map[string]interface{}),
	}
}

// NewDynamicPluginService creates a plugin service with Lua engine for runtime plugin installation
func NewDynamicPluginService(
	commandRegistry *command_handler.CommandHandlerRegistry,
	agentState command_handler.IAgentState,
) *PluginService {
	engine := NewLuaEngine("")
	store := NewPluginStore(engine)

	return &PluginService{
		store:           store,
		commandRegistry: commandRegistry,
		agentState:      agentState,
		config:          make(map[string]interface{}),
	}
}

// GetStore returns the plugin store
func (m *PluginService) GetStore() IPluginStore {
	return m.store
}

// SetConfig sets a configuration value
func (m *PluginService) SetConfig(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config[key] = value
}

// GetConfig gets a configuration value
func (m *PluginService) GetConfig(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.config[key]
	return val, exists
}

// LoadPluginFromFile loads a plugin from a file
func (m *PluginService) LoadPluginFromFile(filePath string) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.store.LoadFromFile(filePath, ctx); err != nil {
		return err
	}

	// Get the loaded plugin
	plugin, _ := m.store.Get(m.getPluginNameFromPath(filePath))
	if plugin != nil {
		if err := m.registerCommandHandler(plugin); err != nil {
			m.store.Unload(plugin.GetMetadata().Name)
			return fmt.Errorf("failed to register command handler: %w", err)
		}
	}

	logging.Info("Plugin loaded from file: %s", filePath)
	return nil
}

// LoadPluginFromBytes loads a plugin from bytes
func (m *PluginService) LoadPluginFromBytes(name string, data []byte) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.store.LoadFromBytes(name, data, ctx); err != nil {
		return err
	}

	// Get the loaded plugin
	plugin, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("plugin '%s' not found after loading", name)
	}

	if err := m.registerCommandHandler(plugin); err != nil {
		m.store.Unload(name)
		return fmt.Errorf("failed to register command handler: %w", err)
	}

	logging.Info("Plugin loaded from bytes: %s", name)
	return nil
}

// UnloadPlugin unloads a plugin
func (m *PluginService) UnloadPlugin(name string) error {
	plugin, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Unregister from command registry
	m.unregisterCommandHandler(plugin)

	if err := m.store.Unload(name); err != nil {
		return err
	}

	logging.Info("Plugin unloaded: %s", name)
	return nil
}

// registerCommandHandler integrates a plugin with the command registry
func (m *PluginService) registerCommandHandler(plugin IPlugin) error {
	metadata := plugin.GetMetadata()

	// Check if it has command handler capability
	hasCapability := false
	for _, cap := range metadata.Capabilities {
		if cap == CapabilityCommandHandler {
			hasCapability = true
			break
		}
	}

	if !hasCapability {
		return nil
	}

	// Check if it's a stateful plugin
	if statefulPlugin, ok := plugin.(IStatefulPlugin); ok {
		handler := command_handler.StatefulCommandHandler{
			Name:    metadata.Name,
			Handler: statefulPlugin.ExecuteStateful,
		}
		m.commandRegistry.RegisterStatefulHandler(handler)
		logging.Debug("Registered stateful command handler: %s", metadata.Name)
	} else {
		// Regular command handler
		handler := command_handler.CommandHandler{
			Name:    metadata.Name,
			Handler: plugin.Execute,
		}
		m.commandRegistry.RegisterHandler(handler)
		logging.Debug("Registered command handler: %s", metadata.Name)
	}

	return nil
}

// unregisterCommandHandler removes a plugin from the command registry
func (m *PluginService) unregisterCommandHandler(plugin IPlugin) {
	metadata := plugin.GetMetadata()
	m.commandRegistry.Unregister(metadata.Name)
	logging.Debug("Unregistered command handler: %s", metadata.Name)
}

// ExecutePlugin executes a plugin command directly
func (m *PluginService) ExecutePlugin(name string, cmd models.Command) (models.CommandResult, error) {
	plugin, exists := m.store.Get(name)
	if !exists {
		return models.CommandResult{}, fmt.Errorf("plugin '%s' not found", name)
	}

	if plugin.GetStatus() != PluginStatusActive {
		return models.CommandResult{}, fmt.Errorf("plugin '%s' is not active", name)
	}

	result := plugin.Execute(cmd)
	return result, nil
}

// ListPlugins returns information about all loaded plugins
func (m *PluginService) ListPlugins() []PluginMetadata {
	plugins := m.store.GetAll()
	metadata := make([]PluginMetadata, 0, len(plugins))

	for _, plugin := range plugins {
		metadata = append(metadata, plugin.GetMetadata())
	}

	return metadata
}

// GetPluginStatus returns the status of a specific plugin
func (m *PluginService) GetPluginStatus(name string) (PluginStatus, error) {
	plugin, exists := m.store.Get(name)
	if !exists {
		return PluginStatusUnloaded, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin.GetStatus(), nil
}

// Shutdown unloads all plugins in reverse order
func (m *PluginService) Shutdown() error {
	order := m.store.LoadOrder()

	// Unload in reverse order to respect dependencies
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		if err := m.UnloadPlugin(name); err != nil {
			logging.Error("Failed to unload plugin '%s' during shutdown: %v", name, err)
		}
	}

	logging.Info("Plugin service shutdown complete")
	return nil
}

// ScanAndLoadPlugins scans the plugin directory and loads all Lua plugins
func (m *PluginService) ScanAndLoadPlugins() error {
	luaEngine, ok := m.store.GetEngine().(*LuaEngine)
	if !ok {
		return fmt.Errorf("ScanAndLoadPlugins only supported for Lua engine")
	}

	if luaEngine.pluginDir == "" {
		return fmt.Errorf("no plugin directory configured")
	}

	plugins, err := scanLuaPluginDirectory(luaEngine.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to scan plugin directory: %w", err)
	}

	for _, pluginPath := range plugins {
		if err := m.LoadPluginFromFile(pluginPath); err != nil {
			logging.Warn("Failed to load plugin %s: %v", pluginPath, err)
		}
	}

	return nil
}

// getPluginNameFromPath extracts plugin name from file path
func (m *PluginService) getPluginNameFromPath(filePath string) string {
	// This is a helper - actual name comes from plugin metadata
	// Just used for initial lookup
	return filePath
}

// scanLuaPluginDirectory scans directory for plugin.lua files in subdirectories
func scanLuaPluginDirectory(pluginDir string) ([]string, error) {
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %w", err)
	}

	plugins := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginFile := filepath.Join(pluginDir, entry.Name(), "plugin.lua")
		if _, err := os.Stat(pluginFile); err == nil {
			plugins = append(plugins, pluginFile)
		}
	}

	return plugins, nil
}
