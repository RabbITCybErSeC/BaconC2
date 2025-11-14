package plugins

import (
	"fmt"
	"sync"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

// PluginManager handles plugin lifecycle and integration with the command system
type PluginManager struct {
	registry        *PluginRegistry
	commandRegistry *command_handler.CommandHandlerRegistry
	agentState      command_handler.IAgentState
	mu              sync.RWMutex
	config          map[string]interface{}
}

// NewPluginManager creates a new plugin manager with filesystem-based loader
func NewPluginManager(
	pluginDir string,
	commandRegistry *command_handler.CommandHandlerRegistry,
	agentState command_handler.IAgentState,
) *PluginManager {
	loader := NewNativePluginLoader(pluginDir)
	registry := NewPluginRegistry(loader)

	return &PluginManager{
		registry:        registry,
		commandRegistry: commandRegistry,
		agentState:      agentState,
		config:          make(map[string]interface{}),
	}
}

// NewDynamicPluginManager creates a plugin manager with dynamic loader for runtime plugin installation
func NewDynamicPluginManager(
	commandRegistry *command_handler.CommandHandlerRegistry,
	agentState command_handler.IAgentState,
) *PluginManager {
	loader := NewDynamicPluginLoader()
	registry := NewPluginRegistry(loader)

	return &PluginManager{
		registry:        registry,
		commandRegistry: commandRegistry,
		agentState:      agentState,
		config:          make(map[string]interface{}),
	}
}

// GetRegistry returns the plugin registry
func (m *PluginManager) GetRegistry() *PluginRegistry {
	return m.registry
}

// SetConfig sets a configuration value
func (m *PluginManager) SetConfig(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config[key] = value
}

// GetConfig gets a configuration value
func (m *PluginManager) GetConfig(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.config[key]
	return val, exists
}

// LoadPluginFromFile loads a plugin from a file
func (m *PluginManager) LoadPluginFromFile(filePath string) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.registry.LoadFromFile(filePath, ctx); err != nil {
		return err
	}

	// Get the loaded plugin
	plugin, _ := m.registry.Get(m.getPluginNameFromPath(filePath))
	if plugin != nil {
		if err := m.registerCommandHandler(plugin); err != nil {
			m.registry.Unload(plugin.GetMetadata().Name)
			return fmt.Errorf("failed to register command handler: %w", err)
		}
	}

	logging.Info("Plugin loaded from file: %s", filePath)
	return nil
}

// LoadPluginFromBytes loads a plugin from bytes
func (m *PluginManager) LoadPluginFromBytes(name string, data []byte) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.registry.LoadFromBytes(name, data, ctx); err != nil {
		return err
	}

	// Get the loaded plugin
	plugin, exists := m.registry.Get(name)
	if !exists {
		return fmt.Errorf("plugin '%s' not found after loading", name)
	}

	if err := m.registerCommandHandler(plugin); err != nil {
		m.registry.Unload(name)
		return fmt.Errorf("failed to register command handler: %w", err)
	}

	logging.Info("Plugin loaded from bytes: %s", name)
	return nil
}

// UnloadPlugin unloads a plugin
func (m *PluginManager) UnloadPlugin(name string) error {
	plugin, exists := m.registry.Get(name)
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Unregister from command registry
	m.unregisterCommandHandler(plugin)

	if err := m.registry.Unload(name); err != nil {
		return err
	}

	logging.Info("Plugin unloaded: %s", name)
	return nil
}

// registerCommandHandler integrates a plugin with the command registry
func (m *PluginManager) registerCommandHandler(plugin IPlugin) error {
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
func (m *PluginManager) unregisterCommandHandler(plugin IPlugin) {
	metadata := plugin.GetMetadata()
	m.commandRegistry.Unregister(metadata.Name)
	logging.Debug("Unregistered command handler: %s", metadata.Name)
}

// ExecutePlugin executes a plugin command directly
func (m *PluginManager) ExecutePlugin(name string, cmd models.Command) (models.CommandResult, error) {
	plugin, exists := m.registry.Get(name)
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
func (m *PluginManager) ListPlugins() []PluginMetadata {
	plugins := m.registry.GetAll()
	metadata := make([]PluginMetadata, 0, len(plugins))

	for _, plugin := range plugins {
		metadata = append(metadata, plugin.GetMetadata())
	}

	return metadata
}

// GetPluginStatus returns the status of a specific plugin
func (m *PluginManager) GetPluginStatus(name string) (PluginStatus, error) {
	plugin, exists := m.registry.Get(name)
	if !exists {
		return PluginStatusUnloaded, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin.GetStatus(), nil
}

// Shutdown unloads all plugins in reverse order
func (m *PluginManager) Shutdown() error {
	order := m.registry.LoadOrder()

	// Unload in reverse order to respect dependencies
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		if err := m.UnloadPlugin(name); err != nil {
			logging.Error("Failed to unload plugin '%s' during shutdown: %v", name, err)
		}
	}

	logging.Info("Plugin manager shutdown complete")
	return nil
}

// ScanAndLoadPlugins scans the plugin directory and loads all plugins
// Only works with filesystem-based loaders (NativePluginLoader)
func (m *PluginManager) ScanAndLoadPlugins() error {
	nativeLoader := m.registry.GetNativeLoader()
	if nativeLoader == nil {
		return fmt.Errorf("ScanAndLoadPlugins only supported for filesystem-based loaders")
	}

	plugins, err := nativeLoader.ScanPluginDirectory()
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

// RegisterPluginDynamic registers a pre-instantiated plugin (for dynamic loaders only)
func (m *PluginManager) RegisterPluginDynamic(name string, plugin IPlugin) error {
	dynamicLoader, ok := m.registry.GetLoader().(*DynamicPluginLoader)
	if !ok {
		return fmt.Errorf("RegisterPluginDynamic only supported for dynamic loaders")
	}

	if err := dynamicLoader.RegisterPlugin(name, plugin); err != nil {
		return err
	}

	// Initialize and register the plugin with the command registry
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := plugin.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %w", name, err)
	}

	// Register command handler
	if err := m.registerCommandHandler(plugin); err != nil {
		return fmt.Errorf("failed to register command handler: %w", err)
	}

	logging.Info("Plugin registered in memory: %s", name)
	return nil
}

// getPluginNameFromPath extracts plugin name from file path
func (m *PluginManager) getPluginNameFromPath(filePath string) string {
	// This is a helper - actual name comes from plugin metadata
	// Just used for initial lookup
	return filePath
}
