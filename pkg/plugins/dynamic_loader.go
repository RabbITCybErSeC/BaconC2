package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
)

// DynamicPluginLoader loads plugins dynamically from bytes with minimal filesystem footprint.
type DynamicPluginLoader struct {
	mu        sync.RWMutex
	instances map[string]IPlugin
	metadata  map[string]PluginMetadata
	tmpDir    string
}

func NewDynamicPluginLoader() *DynamicPluginLoader {
	return &DynamicPluginLoader{
		instances: make(map[string]IPlugin),
		metadata:  make(map[string]PluginMetadata),
		tmpDir:    detectRAMBasedTmpDir(),
	}
}

func NewDynamicPluginLoaderWithTmpDir(tmpDir string) *DynamicPluginLoader {
	return &DynamicPluginLoader{
		instances: make(map[string]IPlugin),
		metadata:  make(map[string]PluginMetadata),
		tmpDir:    tmpDir,
	}
}

func detectRAMBasedTmpDir() string {
	if _, err := os.Stat("/dev/shm"); err == nil {
		return "/dev/shm"
	}
	for _, path := range []string{"/Volumes/RAM", "/private/tmp"} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return os.TempDir()
}

func (l *DynamicPluginLoader) RegisterPlugin(name string, instance IPlugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.instances[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}

	metadata := instance.GetMetadata()
	l.instances[name] = instance
	l.metadata[name] = metadata

	logging.Info("Registered plugin: %s", name)
	return nil
}

func (l *DynamicPluginLoader) GetInstance(name string) (IPlugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	instance, exists := l.instances[name]
	return instance, exists
}

func (l *DynamicPluginLoader) ListLoaded() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.instances))
	for name := range l.instances {
		names = append(names, name)
	}
	return names
}

func (l *DynamicPluginLoader) Unload(name string) error {
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

	logging.Info("Unloaded plugin: %s", name)
	return nil
}

func (l *DynamicPluginLoader) LoadFromFile(filePath string) (IPlugin, error) {
	return nil, fmt.Errorf("LoadFromFile not supported for dynamic loader")
}

func (l *DynamicPluginLoader) LoadFromBytes(name string, data []byte) (IPlugin, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.instances[name]; exists {
		return nil, fmt.Errorf("plugin '%s' already loaded", name)
	}

	tmpFile := filepath.Join(l.tmpDir, fmt.Sprintf("plugin_%s_%d.so", name, os.Getpid()))

	if err := os.WriteFile(tmpFile, data, 0755); err != nil {
		return nil, fmt.Errorf("failed to write temp plugin file: %w", err)
	}

	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			logging.Warn("Failed to remove temp plugin file %s: %v", tmpFile, err)
		}
	}()

	p, err := plugin.Open(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	symNewPlugin, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export 'NewPlugin' function: %w", err)
	}

	newPlugin, ok := symNewPlugin.(func() IPlugin)
	if !ok {
		return nil, fmt.Errorf("NewPlugin has wrong signature")
	}

	instance := newPlugin()
	l.instances[name] = instance
	l.metadata[name] = instance.GetMetadata()

	logging.Info("Loaded plugin from bytes: %s (tmpDir: %s)", name, l.tmpDir)
	return instance, nil
}
