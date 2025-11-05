package main

import (
	"fmt"
	"strings"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

// EchoPlugin is a simple plugin that echoes back command arguments
type EchoPlugin struct {
	metadata plugins.PluginMetadata
	status   plugins.PluginStatus
	context  *plugins.PluginContext
}

// NewPlugin is the required factory function that must be exported
func NewPlugin() plugins.IPlugin {
	return &EchoPlugin{
		metadata: plugins.PluginMetadata{
			Name:        "echo",
			Version:     "1.0.0",
			Author:      "BaconC2",
			Description: "Echoes back the provided arguments",
			Capabilities: []plugins.PluginCapability{
				plugins.CapabilityCommandHandler,
			},
			Dependencies: []string{},
			Status:       plugins.PluginStatusUnloaded,
		},
		status: plugins.PluginStatusUnloaded,
	}
}

// GetMetadata returns plugin metadata
func (p *EchoPlugin) GetMetadata() plugins.PluginMetadata {
	return p.metadata
}

// Initialize sets up the plugin
func (p *EchoPlugin) Initialize(ctx *plugins.PluginContext) error {
	p.context = ctx
	p.status = plugins.PluginStatusActive
	p.metadata.Status = plugins.PluginStatusActive
	return nil
}

// Execute handles the echo command
func (p *EchoPlugin) Execute(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusCompleted,
			Output:     "Echo: (empty)",
			ResultType: models.ResultTypeText,
		}
	}

	output := fmt.Sprintf("Echo: %s", strings.Join(cmd.Args, " "))

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     output,
		ResultType: models.ResultTypeText,
	}
}

// Cleanup performs cleanup operations
func (p *EchoPlugin) Cleanup() error {
	p.status = plugins.PluginStatusUnloaded
	p.metadata.Status = plugins.PluginStatusUnloaded
	p.context = nil
	return nil
}

// GetStatus returns the current plugin status
func (p *EchoPlugin) GetStatus() plugins.PluginStatus {
	return p.status
}
