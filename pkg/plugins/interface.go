package plugins

import (
	"time"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

// PluginCapability represents what a plugin can do
type PluginCapability string

const (
	CapabilityCommandHandler PluginCapability = "command_handler"
	CapabilityTransport      PluginCapability = "transport"
	CapabilityEncoder        PluginCapability = "encoder"
	CapabilityPersistence    PluginCapability = "persistence"
	CapabilityEvasion        PluginCapability = "evasion"
	CapabilityCustom         PluginCapability = "custom"
)

// PluginStatus represents the current state of a plugin
type PluginStatus string

const (
	PluginStatusUnloaded  PluginStatus = "unloaded"
	PluginStatusLoading   PluginStatus = "loading"
	PluginStatusLoaded    PluginStatus = "loaded"
	PluginStatusActive    PluginStatus = "active"
	PluginStatusError     PluginStatus = "error"
	PluginStatusUnloading PluginStatus = "unloading"
)

// PluginMetadata contains information about a plugin
type PluginMetadata struct {
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Author       string             `json:"author"`
	Description  string             `json:"description"`
	Capabilities []PluginCapability `json:"capabilities"`
	Dependencies []string           `json:"dependencies"`
	LoadedAt     time.Time          `json:"loaded_at"`
	Status       PluginStatus       `json:"status"`
	FilePath     string             `json:"file_path,omitempty"`
	Hash         string             `json:"hash,omitempty"`
}

// PluginContext provides runtime context to plugins
type PluginContext struct {
	AgentState   command_handler.IAgentState
	Registry     *command_handler.CommandHandlerRegistry
	Config       map[string]interface{}
}

// IPlugin is the main interface that all plugins must implement
// This is the interface that compiled .so files must export
type IPlugin interface {
	// GetMetadata returns plugin metadata
	GetMetadata() PluginMetadata

	// Initialize is called when the plugin is loaded
	Initialize(ctx *PluginContext) error

	// Execute handles a command and returns a result
	Execute(cmd models.Command) models.CommandResult

	// Cleanup is called when the plugin is unloaded
	Cleanup() error

	// GetStatus returns the current status of the plugin
	GetStatus() PluginStatus
}

// IStatefulPlugin extends IPlugin with stateful command handling
type IStatefulPlugin interface {
	IPlugin

	// ExecuteStateful handles a command with access to agent state
	ExecuteStateful(ctx *command_handler.CommandContext) models.CommandResult
}

// PluginFactory is the function signature that plugin .so files must export
// The symbol name must be "NewPlugin"
type PluginFactory func() IPlugin

// PluginInfo contains information about a plugin file
type PluginInfo struct {
	Metadata PluginMetadata `json:"metadata"`
	FilePath string         `json:"file_path"`
	FileSize int64          `json:"file_size"`
	Hash     string         `json:"hash"`
}
