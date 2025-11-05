package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

// PluginCompiler handles server-side plugin compilation
type PluginCompiler struct {
	sourceDir  string
	outputDir  string
	buildFlags []string
}

// NewPluginCompiler creates a new plugin compiler
func NewPluginCompiler(sourceDir, outputDir string) *PluginCompiler {
	return &PluginCompiler{
		sourceDir:  sourceDir,
		outputDir:  outputDir,
		buildFlags: []string{},
	}
}

// SetBuildFlags sets custom build flags
func (c *PluginCompiler) SetBuildFlags(flags []string) {
	c.buildFlags = flags
}

// CompilePlugin compiles a Go plugin from source
func (c *PluginCompiler) CompilePlugin(pluginName string, sourceFiles []string) (*plugins.PluginInfo, error) {
	// Validate OS - plugins only work on Linux and macOS
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("Go plugins are not supported on Windows")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build output path
	outputPath := filepath.Join(c.outputDir, pluginName+".so")

	// Prepare build command
	args := []string{"build", "-buildmode=plugin"}
	args = append(args, c.buildFlags...)
	args = append(args, "-o", outputPath)
	args = append(args, sourceFiles...)

	// Execute build
	cmd := exec.Command("go", args...)
	cmd.Dir = c.sourceDir
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w\nOutput: %s", err, string(output))
	}

	logging.Info("Plugin compiled successfully: %s", outputPath)

	// Get file info
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat compiled plugin: %w", err)
	}

	// Calculate hash
	hash, err := plugins.CalculateHash(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	pluginInfo := &plugins.PluginInfo{
		FilePath: outputPath,
		FileSize: fileInfo.Size(),
		Hash:     hash,
	}

	return pluginInfo, nil
}

// CompilePluginFromDirectory compiles all .go files in a directory
func (c *PluginCompiler) CompilePluginFromDirectory(pluginName, sourceDir string) (*plugins.PluginInfo, error) {
	// Find all .go files in the directory
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}

	sourceFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			sourceFiles = append(sourceFiles, filepath.Join(sourceDir, entry.Name()))
		}
	}

	if len(sourceFiles) == 0 {
		return nil, fmt.Errorf("no Go source files found in %s", sourceDir)
	}

	return c.CompilePlugin(pluginName, sourceFiles)
}

// ReadCompiledPlugin reads a compiled plugin file
func (c *PluginCompiler) ReadCompiledPlugin(pluginName string) ([]byte, error) {
	pluginPath := filepath.Join(c.outputDir, pluginName+".so")
	
	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin file: %w", err)
	}

	return data, nil
}

// ListCompiledPlugins lists all compiled plugins in the output directory
func (c *PluginCompiler) ListCompiledPlugins() ([]string, error) {
	if _, err := os.Stat(c.outputDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(c.outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read output directory: %w", err)
	}

	plugins := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".so" {
			plugins = append(plugins, strings.TrimSuffix(entry.Name(), ".so"))
		}
	}

	return plugins, nil
}

// DeleteCompiledPlugin deletes a compiled plugin
func (c *PluginCompiler) DeleteCompiledPlugin(pluginName string) error {
	pluginPath := filepath.Join(c.outputDir, pluginName+".so")
	
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}

	logging.Info("Deleted compiled plugin: %s", pluginPath)
	return nil
}

// GetPluginPath returns the full path to a compiled plugin
func (c *PluginCompiler) GetPluginPath(pluginName string) string {
	return filepath.Join(c.outputDir, pluginName+".so")
}

// ValidatePluginSource performs basic validation on plugin source code
func (c *PluginCompiler) ValidatePluginSource(sourceFiles []string) error {
	for _, file := range sourceFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("source file not found: %s", file)
		}
	}
	return nil
}
