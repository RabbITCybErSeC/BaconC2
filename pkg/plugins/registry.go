package plugins

import (
	"fmt"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
)

// PluginRegistry manages all loaded plugins in memory
type PluginRegistry struct {
	mu      sync.RWMutex
	plugins map[string]*PluginEntry
	order   []string // Track load order for dependency management
	loader  *NativePluginLoader
}

// PluginEntry wraps a plugin with additional metadata
type PluginEntry struct {
	Plugin   IPlugin
	FilePath string
	LoadedAt time.Time
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry(loader *NativePluginLoader) *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]*PluginEntry),
		order:   make([]string, 0),
		loader:  loader,
	}
}

// LoadFromFile loads a plugin from a file path
func (r *PluginRegistry) LoadFromFile(filePath string, ctx *PluginContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load the plugin using the native loader
	plugin, err := r.loader.LoadFromFile(filePath)
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
func (r *PluginRegistry) LoadFromBytes(name string, data []byte, ctx *PluginContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load the plugin using the native loader
	plugin, err := r.loader.LoadFromBytes(name, data)
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
		FilePath: r.loader.GetPluginDir() + "/" + name + ".so",
		LoadedAt: time.Now(),
	}

	r.plugins[metadata.Name] = entry
	r.order = append(r.order, metadata.Name)

	logging.Info("Plugin '%s' loaded successfully from bytes", metadata.Name)
	return nil
}

// Unload deactivates and removes a plugin
func (r *PluginRegistry) Unload(name string) error {
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

	// Unload from native loader
	if err := r.loader.Unload(entry.FilePath); err != nil {
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
func (r *PluginRegistry) Get(name string) (IPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.plugins[name]
	if !exists {
		return nil, false
	}
	return entry.Plugin, true
}

// GetEntry retrieves a plugin entry by name
func (r *PluginRegistry) GetEntry(name string) (*PluginEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.plugins[name]
	return entry, exists
}

// GetAll returns all loaded plugins
func (r *PluginRegistry) GetAll() map[string]IPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make(map[string]IPlugin, len(r.plugins))
	for name, entry := range r.plugins {
		plugins[name] = entry.Plugin
	}
	return plugins
}

// GetAllEntries returns all plugin entries
func (r *PluginRegistry) GetAllEntries() map[string]*PluginEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make(map[string]*PluginEntry, len(r.plugins))
	for name, entry := range r.plugins {
		entries[name] = entry
	}
	return entries
}

// IsLoaded checks if a plugin is currently loaded
func (r *PluginRegistry) IsLoaded(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.plugins[name]
	return exists
}

// Count returns the total number of loaded plugins
func (r *PluginRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}

// LoadOrder returns the order in which plugins were loaded
func (r *PluginRegistry) LoadOrder() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order := make([]string, len(r.order))
	copy(order, r.order)
	return order
}

// GetByCapability returns all plugins with a specific capability
func (r *PluginRegistry) GetByCapability(capability PluginCapability) []IPlugin {
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

// GetLoader returns the native plugin loader
func (r *PluginRegistry) GetLoader() *NativePluginLoader {
	return r.loader
}
