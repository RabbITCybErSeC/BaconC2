package plugins

import (
	"time"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type PluginCapability string

const (
	CapabilityCommandHandler PluginCapability = "command_handler"
	CapabilityTransport      PluginCapability = "transport"
	CapabilityEncoder        PluginCapability = "encoder"
	CapabilityPersistence    PluginCapability = "persistence"
	CapabilityEvasion        PluginCapability = "evasion"
	CapabilityCustom         PluginCapability = "custom"
)

type PluginStatus string

const (
	PluginStatusUnloaded  PluginStatus = "unloaded"
	PluginStatusLoading   PluginStatus = "loading"
	PluginStatusLoaded    PluginStatus = "loaded"
	PluginStatusActive    PluginStatus = "active"
	PluginStatusError     PluginStatus = "error"
	PluginStatusUnloading PluginStatus = "unloading"
)

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

type PluginContext struct {
	AgentState command_handler.IAgentState
	Registry   *command_handler.CommandHandlerRegistry
	Config     map[string]interface{}
}

type IPlugin interface {
	GetMetadata() PluginMetadata
	Initialize(ctx *PluginContext) error
	Execute(cmd models.Command) models.CommandResult
	Cleanup() error
	GetStatus() PluginStatus
}

type IStatefulPlugin interface {
	IPlugin
	ExecuteStateful(ctx *command_handler.CommandContext) models.CommandResult
}

type PluginFactory func() IPlugin

type PluginInfo struct {
	Metadata PluginMetadata `json:"metadata"`
	FilePath string         `json:"file_path"`
	FileSize int64          `json:"file_size"`
	Hash     string         `json:"hash"`
}

// PluginEntry wraps a plugin with additional metadata
type PluginEntry struct {
	Plugin   IPlugin
	FilePath string
	LoadedAt time.Time
}

// IPluginEngine defines the interface for plugin execution engines
// An engine is responsible for loading, executing, and managing plugin instances
type IPluginEngine interface {
	GetInstance(identifier string) (IPlugin, bool)
	ListLoaded() []string
	Unload(identifier string) error
	LoadFromFile(filePath string) (IPlugin, error)
	LoadFromBytes(name string, data []byte) (IPlugin, error)
}

// IPluginStore defines the interface for plugin storage and management
// A store manages loaded plugins, tracks dependencies, and provides query capabilities
type IPluginStore interface {
	LoadFromFile(filePath string, ctx *PluginContext) error
	LoadFromBytes(name string, data []byte, ctx *PluginContext) error
	Unload(name string) error
	Get(name string) (IPlugin, bool)
	GetEntry(name string) (*PluginEntry, bool)
	GetAll() map[string]IPlugin
	GetAllEntries() map[string]*PluginEntry
	IsLoaded(name string) bool
	Count() int
	LoadOrder() []string
	GetByCapability(capability PluginCapability) []IPlugin
	GetEngine() IPluginEngine
	GetLuaEngine() *LuaEngine
}
