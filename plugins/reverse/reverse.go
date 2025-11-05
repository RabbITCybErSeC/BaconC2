package main

import (
	"strings"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
)

// ReversePlugin reverses string arguments
type ReversePlugin struct {
	metadata plugins.PluginMetadata
	status   plugins.PluginStatus
	context  *plugins.PluginContext
}

// NewPlugin is the required factory function that must be exported
func NewPlugin() plugins.IPlugin {
	return &ReversePlugin{
		metadata: plugins.PluginMetadata{
			Name:        "reverse",
			Version:     "1.0.0",
			Author:      "BaconC2",
			Description: "Reverses the provided string arguments",
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
func (p *ReversePlugin) GetMetadata() plugins.PluginMetadata {
	return p.metadata
}

// Initialize sets up the plugin
func (p *ReversePlugin) Initialize(ctx *plugins.PluginContext) error {
	p.context = ctx
	p.status = plugins.PluginStatusActive
	p.metadata.Status = plugins.PluginStatusActive
	return nil
}

// Execute handles the reverse command
func (p *ReversePlugin) Execute(cmd models.Command) models.CommandResult {
	if len(cmd.Args) == 0 {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Error: No arguments provided",
			ResultType: models.ResultTypeError,
		}
	}

	input := strings.Join(cmd.Args, " ")
	reversed := reverseString(input)

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     reversed,
		ResultType: models.ResultTypeText,
	}
}

// Cleanup performs cleanup operations
func (p *ReversePlugin) Cleanup() error {
	p.status = plugins.PluginStatusUnloaded
	p.metadata.Status = plugins.PluginStatusUnloaded
	p.context = nil
	return nil
}

// GetStatus returns the current plugin status
func (p *ReversePlugin) GetStatus() plugins.PluginStatus {
	return p.status
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
