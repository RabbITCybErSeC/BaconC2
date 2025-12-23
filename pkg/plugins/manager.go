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

type PluginService struct {
	store           IPluginStore
	commandRegistry *command_handler.CommandHandlerRegistry
	agentState      command_handler.IAgentState
	mu              sync.RWMutex
	config          map[string]interface{}
}

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

func (m *PluginService) GetStore() IPluginStore {
	return m.store
}

func (m *PluginService) SetConfig(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config[key] = value
}

func (m *PluginService) GetConfig(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.config[key]
	return val, exists
}

func (m *PluginService) LoadPluginFromFile(filePath string) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.store.LoadFromFile(filePath, ctx); err != nil {
		return err
	}

	order := m.store.LoadOrder()
	if len(order) == 0 {
		return fmt.Errorf("no plugins loaded")
	}

	lastPluginName := order[len(order)-1]
	plugin, exists := m.store.Get(lastPluginName)
	if !exists {
		return fmt.Errorf("plugin not found after loading")
	}

	if err := m.registerCommandHandler(plugin); err != nil {
		m.store.Unload(lastPluginName)
		return fmt.Errorf("failed to register command handler: %w", err)
	}

	logging.Info("Plugin loaded from file: %s", filePath)
	return nil
}

func (m *PluginService) LoadPluginFromBytes(name string, data []byte) error {
	ctx := &PluginContext{
		AgentState: m.agentState,
		Registry:   m.commandRegistry,
		Config:     m.config,
	}

	if err := m.store.LoadFromBytes(name, data, ctx); err != nil {
		return err
	}

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

func (m *PluginService) UnloadPlugin(name string) error {
	plugin, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	m.unregisterCommandHandler(plugin)

	if err := m.store.Unload(name); err != nil {
		return err
	}

	logging.Info("Plugin unloaded: %s", name)
	return nil
}

func (m *PluginService) registerCommandHandler(plugin IPlugin) error {
	metadata := plugin.GetMetadata()

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

	if statefulPlugin, ok := plugin.(IStatefulPlugin); ok {
		handler := command_handler.StatefulCommandHandler{
			Name:    metadata.Name,
			Handler: statefulPlugin.ExecuteStateful,
		}
		m.commandRegistry.RegisterStatefulHandler(handler)
		logging.Debug("Registered stateful command handler: %s", metadata.Name)
	} else {
		handler := command_handler.CommandHandler{
			Name:    metadata.Name,
			Handler: plugin.Execute,
		}
		m.commandRegistry.RegisterHandler(handler)
		logging.Debug("Registered command handler: %s", metadata.Name)
	}

	return nil
}

func (m *PluginService) unregisterCommandHandler(plugin IPlugin) {
	metadata := plugin.GetMetadata()
	m.commandRegistry.Unregister(metadata.Name)
	logging.Debug("Unregistered command handler: %s", metadata.Name)
}

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

func (m *PluginService) ListPlugins() []PluginMetadata {
	plugins := m.store.GetAll()
	metadata := make([]PluginMetadata, 0, len(plugins))

	for _, plugin := range plugins {
		metadata = append(metadata, plugin.GetMetadata())
	}

	return metadata
}

func (m *PluginService) GetPluginStatus(name string) (PluginStatus, error) {
	plugin, exists := m.store.Get(name)
	if !exists {
		return PluginStatusUnloaded, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin.GetStatus(), nil
}

func (m *PluginService) Shutdown() error {
	order := m.store.LoadOrder()

	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		if err := m.UnloadPlugin(name); err != nil {
			logging.Error("Failed to unload plugin '%s' during shutdown: %v", name, err)
		}
	}

	logging.Info("Plugin service shutdown complete")
	return nil
}

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
