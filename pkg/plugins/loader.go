package plugins

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
)

// NativePluginLoader handles loading and unloading of native Go plugins (.so files)
type NativePluginLoader struct {
	mu             sync.RWMutex
	loadedPlugins  map[string]*plugin.Plugin
	pluginDir      string
	instances      map[string]IPlugin
}

// NewNativePluginLoader creates a new plugin loader
func NewNativePluginLoader(pluginDir string) *NativePluginLoader {
	return &NativePluginLoader{
		loadedPlugins: make(map[string]*plugin.Plugin),
		instances:     make(map[string]IPlugin),
		pluginDir:     pluginDir,
	}
}

// LoadFromFile loads a plugin from a .so file
func (l *NativePluginLoader) LoadFromFile(filePath string) (IPlugin, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if already loaded
	if instance, exists := l.instances[filePath]; exists {
		return instance, nil
	}

	// Load the plugin
	p, err := plugin.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the NewPlugin symbol
	symNewPlugin, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export 'NewPlugin' function: %w", err)
	}

	// Assert the symbol is a PluginFactory
	newPlugin, ok := symNewPlugin.(func() IPlugin)
	if !ok {
		return nil, fmt.Errorf("NewPlugin has wrong signature")
	}

	// Create plugin instance
	instance := newPlugin()
	
	l.loadedPlugins[filePath] = p
	l.instances[filePath] = instance

	logging.Info("Loaded plugin from: %s", filePath)
	return instance, nil
}

// LoadFromBytes saves plugin bytes to disk and loads it
func (l *NativePluginLoader) LoadFromBytes(name string, data []byte) (IPlugin, error) {
	// Ensure plugin directory exists
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Generate file path
	filePath := filepath.Join(l.pluginDir, name+".so")

	// Write plugin to disk
	if err := os.WriteFile(filePath, data, 0755); err != nil {
		return nil, fmt.Errorf("failed to write plugin file: %w", err)
	}

	logging.Info("Saved plugin to: %s", filePath)

	// Load the plugin
	return l.LoadFromFile(filePath)
}

// Unload removes a plugin from memory
// Note: Go's plugin system doesn't support true unloading, but we can remove our references
func (l *NativePluginLoader) Unload(filePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	instance, exists := l.instances[filePath]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", filePath)
	}

	// Call cleanup on the plugin
	if err := instance.Cleanup(); err != nil {
		logging.Warn("Plugin cleanup error: %v", err)
	}

	// Remove references (note: the .so remains in memory due to Go limitations)
	delete(l.instances, filePath)
	delete(l.loadedPlugins, filePath)

	logging.Info("Unloaded plugin: %s", filePath)
	return nil
}

// GetInstance returns a loaded plugin instance
func (l *NativePluginLoader) GetInstance(filePath string) (IPlugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	instance, exists := l.instances[filePath]
	return instance, exists
}

// ListLoaded returns all loaded plugin file paths
func (l *NativePluginLoader) ListLoaded() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	paths := make([]string, 0, len(l.instances))
	for path := range l.instances {
		paths = append(paths, path)
	}
	return paths
}

// GetPluginDir returns the plugin directory
func (l *NativePluginLoader) GetPluginDir() string {
	return l.pluginDir
}

// CalculateHash computes SHA256 hash of a file
func CalculateHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// CalculateHashBytes computes SHA256 hash of bytes
func CalculateHashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:])
}

// ScanPluginDirectory scans the plugin directory for .so files
func (l *NativePluginLoader) ScanPluginDirectory() ([]string, error) {
	if _, err := os.Stat(l.pluginDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(l.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %w", err)
	}

	plugins := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".so" {
			plugins = append(plugins, filepath.Join(l.pluginDir, entry.Name()))
		}
	}

	return plugins, nil
}

// DeletePlugin removes a plugin file from disk
func (l *NativePluginLoader) DeletePlugin(filePath string) error {
	// First unload if loaded
	if _, exists := l.instances[filePath]; exists {
		if err := l.Unload(filePath); err != nil {
			return err
		}
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete plugin file: %w", err)
	}

	logging.Info("Deleted plugin file: %s", filePath)
	return nil
}
