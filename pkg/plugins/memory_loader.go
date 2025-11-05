package plugins

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
)

// InMemoryPluginLoader handles loading and unloading of plugins entirely in memory
// without touching the filesystem. It implements IPluginLoader interface.
type InMemoryPluginLoader struct {
	mu        sync.RWMutex
	instances map[string]IPlugin
	metadata  map[string]PluginMetadata
}

// NewInMemoryPluginLoader creates a new in-memory plugin loader
func NewInMemoryPluginLoader() *InMemoryPluginLoader {
	return &InMemoryPluginLoader{
		instances: make(map[string]IPlugin),
		metadata:  make(map[string]PluginMetadata),
	}
}

// RegisterPlugin registers a pre-instantiated plugin in memory
func (l *InMemoryPluginLoader) RegisterPlugin(name string, instance IPlugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.instances[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}

	metadata := instance.GetMetadata()
	l.instances[name] = instance
	l.metadata[name] = metadata

	logging.Info("Registered plugin in memory: %s", name)
	return nil
}

// GetInstance returns a loaded plugin instance
func (l *InMemoryPluginLoader) GetInstance(name string) (IPlugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	instance, exists := l.instances[name]
	return instance, exists
}

// ListLoaded returns all loaded plugin names
func (l *InMemoryPluginLoader) ListLoaded() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.instances))
	for name := range l.instances {
		names = append(names, name)
	}
	return names
}

// Unload removes a plugin from memory
func (l *InMemoryPluginLoader) Unload(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	instance, exists := l.instances[name]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", name)
	}

	if err := instance.Cleanup(); err != nil {
		logging.Warn("Plugin cleanup error: %v", err)
	}

	delete(l.instances, name)
	delete(l.metadata, name)

	logging.Info("Unloaded plugin from memory: %s", name)
	return nil
}

// LoadFromFile is not supported for in-memory loader
func (l *InMemoryPluginLoader) LoadFromFile(filePath string) (IPlugin, error) {
	return nil, fmt.Errorf("LoadFromFile not supported for in-memory loader")
}

// LoadFromBytes is not supported for in-memory loader
func (l *InMemoryPluginLoader) LoadFromBytes(name string, data []byte) (IPlugin, error) {
	return nil, fmt.Errorf("LoadFromBytes not supported for in-memory loader")
}

// CalculateHashBytes computes SHA256 hash of bytes
func CalculateHashBytesInMemory(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:])
}
