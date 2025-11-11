package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <plugin.so> <output.json>\n", os.Args[0])
		os.Exit(1)
	}

	pluginPath := os.Args[1]
	outputPath := os.Args[2]

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load plugin: %v\n", err)
		os.Exit(1)
	}

	// Look up the NewPlugin symbol
	symNewPlugin, err := p.Lookup("NewPlugin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find NewPlugin function: %v\n", err)
		os.Exit(1)
	}

	// Cast to the correct function type
	newPlugin, ok := symNewPlugin.(func() plugins.IPlugin)
	if !ok {
		fmt.Fprintf(os.Stderr, "NewPlugin has wrong signature\n")
		os.Exit(1)
	}

	// Create plugin instance
	pluginInstance := newPlugin()

	// Get metadata
	metadata := pluginInstance.GetMetadata()

	// Add file information
	fileInfo, err := os.Stat(pluginPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stat plugin file: %v\n", err)
		os.Exit(1)
	}

	pluginInfo := plugins.PluginInfo{
		Metadata: metadata,
		FilePath: filepath.Base(pluginPath),
		FileSize: fileInfo.Size(),
		Hash:     "", // Will be calculated by the server
	}

	// Write to JSON
	data, err := json.MarshalIndent(pluginInfo, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Metadata extracted to %s\n", outputPath)
}
