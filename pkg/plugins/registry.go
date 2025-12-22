package plugins

import (
	"fmt"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
)

// PluginStore manages all loaded plugins in memory
// It tracks plugin entries, load order, and dependencies
type PluginStore struct {
	mu      sync.RWMutex
	plugins map[string]*PluginEntry
	order   []string // Track load order for dependency management
	engine  IPluginEngine
}

// PluginEntry wraps a plugin with additional metadata
type PluginEntry struct {
	Plugin   IPlugin
	FilePath string
	LoadedAt time.Time
}

// NewPluginStore creates a new plugin store
func NewPluginStore(engine IPluginEngine) *PluginStore {
	return &PluginStore{
		plugins: make(map[string]*PluginEntry),
		order:   make([]string, 0),
		engine:  engine,
	}
}

// LoadFromFile loads a plugin from a file path
func (r *PluginStore) LoadFromFile(filePath string, ctx *PluginContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load the plugin using the engine
	plugin, err := r.engine.LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}

	metadata := plugin.GetMetadata()

	// Check if already registered
	if _, exists := r.plugins[metadata.Name]; exists {
		return fmt.Errorf("plugin '%s' already loaded", metadata.Name)
	}

	// Check dependencies
	for _, dep := range metadata.Dependencies {
		if _, exists := r.plugins[dep]; !exists {
			return fmt.Errorf("dependency '%s' not loaded for plugin '%s'", dep, metadata.Name)
		}
	}

	// Initialize the plugin
	if err := plugin.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %w", metadata.Name, err)
	}

	// Register the plugin
	entry := &PluginEntry{
		Plugin:   plugin,
		FilePath: filePath,
		LoadedAt: time.Now(),
	}

	r.plugins[metadata.Name] = entry
	r.order = append(r.order, metadata.Name)

	logging.Info("Plugin '%s' loaded successfully from %s", metadata.Name, filePath)
	return nil
}

// LoadFromBytes loads a plugin from bytes
func (r *PluginStore) LoadFromBytes(name string, data []byte, ctx *PluginContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load the plugin using the engine
	plugin, err := r.engine.LoadFromBytes(name, data)
	if err != nil {
		return fmt.Errorf("failed to load plugin from bytes: %w", err)
	}

	metadata := plugin.GetMetadata()

	// Check if already registered
	if _, exists := r.plugins[metadata.Name]; exists {
		return fmt.Errorf("plugin '%s' already loaded", metadata.Name)
	}

	// Check dependencies
	for _, dep := range metadata.Dependencies {
		if _, exists := r.plugins[dep]; !exists {
			return fmt.Errorf("dependency '%s' not loaded for plugin '%s'", dep, metadata.Name)
		}
	}

	// Initialize the plugin
	if err := plugin.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %w", metadata.Name, err)
	}

	// Register the plugin
	entry := &PluginEntry{
		Plugin:   plugin,
		FilePath: name,
		LoadedAt: time.Now(),
	}

	r.plugins[metadata.Name] = entry
	r.order = append(r.order, metadata.Name)

	logging.Info("Plugin '%s' loaded successfully from bytes", metadata.Name)
	return nil
}

// Unload deactivates and removes a plugin
func (r *PluginStore) Unload(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Check if any loaded plugins depend on this one
	for _, otherName := range r.order {
		if otherName == name {
			continue
		}
		otherEntry, ok := r.plugins[otherName]
		if !ok {
			continue
		}
		otherMetadata := otherEntry.Plugin.GetMetadata()
		for _, dep := range otherMetadata.Dependencies {
			if dep == name {
				return fmt.Errorf("cannot unload '%s': plugin '%s' depends on it", name, otherName)
			}
		}
	}

	// Unload from engine
	if err := r.engine.Unload(entry.FilePath); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	// Remove from registry
	delete(r.plugins, name)

	// Remove from load order
	for i, n := range r.order {
		if n == name {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}

	logging.Info("Plugin '%s' unloaded successfully", name)
	return nil
}

// Get retrieves a plugin by name
func (r *PluginStore) Get(name string) (IPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.plugins[name]
	if !exists {
		return nil, false
	}
	return entry.Plugin, true
}

// GetEntry retrieves a plugin entry by name
func (r *PluginStore) GetEntry(name string) (*PluginEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.plugins[name]
	return entry, exists
}

// GetAll returns all loaded plugins
func (r *PluginStore) GetAll() map[string]IPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make(map[string]IPlugin, len(r.plugins))
	for name, entry := range r.plugins {
		plugins[name] = entry.Plugin
	}
	return plugins
}

// GetAllEntries returns all plugin entries
func (r *PluginStore) GetAllEntries() map[string]*PluginEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make(map[string]*PluginEntry, len(r.plugins))
	for name, entry := range r.plugins {
		entries[name] = entry
	}
	return entries
}

// IsLoaded checks if a plugin is currently loaded
func (r *PluginStore) IsLoaded(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.plugins[name]
	return exists
}

// Count returns the total number of loaded plugins
func (r *PluginStore) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}

// LoadOrder returns the order in which plugins were loaded
func (r *PluginStore) LoadOrder() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order := make([]string, len(r.order))
	copy(order, r.order)
	return order
}

// GetByCapability returns all plugins with a specific capability
func (r *PluginStore) GetByCapability(capability PluginCapability) []IPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]IPlugin, 0)
	for _, entry := range r.plugins {
		metadata := entry.Plugin.GetMetadata()
		for _, cap := range metadata.Capabilities {
			if cap == capability {
				plugins = append(plugins, entry.Plugin)
				break
			}
		}
	}
	return plugins
}

// GetEngine returns the plugin engine
func (r *PluginStore) GetEngine() IPluginEngine {
	return r.engine
}

// GetLuaEngine returns the Lua engine if available
func (r *PluginStore) GetLuaEngine() *LuaEngine {
	if luaEngine, ok := r.engine.(*LuaEngine); ok {
		return luaEngine
	}
	return nil
}
