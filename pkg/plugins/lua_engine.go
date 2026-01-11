package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	lua "github.com/yuin/gopher-lua"
)

type LuaEngine struct {
	mu        sync.RWMutex
	instances map[string]*LuaPlugin
	pluginDir string
}

type LuaPlugin struct {
	name     string
	vm       *lua.LState
	metadata PluginMetadata
	status   PluginStatus
	context  *PluginContext
}

func NewLuaEngine(pluginDir string) IPluginEngine {
	return &LuaEngine{
		instances: make(map[string]*LuaPlugin),
		pluginDir: pluginDir,
	}
}

func (l *LuaEngine) LoadFromFile(filePath string) (IPlugin, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin file: %w", err)
	}

	name := filepath.Base(filePath)
	name = name[:len(name)-len(filepath.Ext(name))]

	return l.loadFromBytesUnlocked(name, data)
}

func (l *LuaEngine) LoadFromBytes(name string, data []byte) (IPlugin, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.loadFromBytesUnlocked(name, data)
}

func (l *LuaEngine) loadFromBytesUnlocked(name string, data []byte) (IPlugin, error) {

	vm := lua.NewState()

	if err := vm.DoString(string(data)); err != nil {
		vm.Close()
		return nil, fmt.Errorf("failed to execute plugin script: %w", err)
	}

	metadataTable := vm.GetGlobal("metadata")
	if metadataTable == lua.LNil {
		vm.Close()
		return nil, fmt.Errorf("plugin missing 'metadata' table")
	}

	table, ok := metadataTable.(*lua.LTable)
	if !ok {
		vm.Close()
		return nil, fmt.Errorf("metadata must be a table, got %s", metadataTable.Type())
	}

	metadata, err := l.parseMetadata(vm, table)
	if err != nil {
		vm.Close()
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	if metadata.Name == "" {
		metadata.Name = name
	}
	metadata.Status = PluginStatusUnloaded
	metadata.LoadedAt = time.Now()

	plugin := &LuaPlugin{
		name:     metadata.Name,
		vm:       vm,
		metadata: metadata,
		status:   PluginStatusUnloaded,
	}

	if _, exists := l.instances[metadata.Name]; exists {
		return nil, fmt.Errorf("plugin '%s' already loaded", metadata.Name)
	}

	l.instances[metadata.Name] = plugin
	logging.Info("Loaded Lua plugin: %s", metadata.Name)
	return plugin, nil
}

func (l *LuaEngine) parseMetadata(vm *lua.LState, table *lua.LTable) (PluginMetadata, error) {
	metadata := PluginMetadata{
		Capabilities: []PluginCapability{},
		Dependencies: []string{},
	}

	table.ForEach(func(key, value lua.LValue) {
		keyStr := key.String()
		switch keyStr {
		case "name":
			metadata.Name = value.String()
		case "version":
			metadata.Version = value.String()
		case "author":
			metadata.Author = value.String()
		case "description":
			metadata.Description = value.String()
		case "capabilities":
			if capTable, ok := value.(*lua.LTable); ok {
				capTable.ForEach(func(_, cap lua.LValue) {
					metadata.Capabilities = append(metadata.Capabilities, PluginCapability(cap.String()))
				})
			}
		case "dependencies":
			if depTable, ok := value.(*lua.LTable); ok {
				depTable.ForEach(func(_, dep lua.LValue) {
					metadata.Dependencies = append(metadata.Dependencies, dep.String())
				})
			}
		}
	})

	return metadata, nil
}

func (l *LuaEngine) Unload(identifier string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	plugin, exists := l.instances[identifier]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", identifier)
	}

	if err := plugin.Cleanup(); err != nil {
		logging.Warn("Plugin cleanup error: %v", err)
	}

	plugin.vm.Close()
	delete(l.instances, identifier)

	logging.Info("Unloaded Lua plugin: %s", identifier)
	return nil
}

func (l *LuaEngine) GetInstance(identifier string) (IPlugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	instance, exists := l.instances[identifier]
	return instance, exists
}

func (l *LuaEngine) ListLoaded() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.instances))
	for name := range l.instances {
		names = append(names, name)
	}
	return names
}

func (p *LuaPlugin) GetMetadata() PluginMetadata {
	return p.metadata
}

func (p *LuaPlugin) Initialize(ctx *PluginContext) error {
	p.context = ctx

	// Expose context objects as Lua tables with serialized JSON data
	if ctx.AgentState != nil {
		if agentStateJSON, err := json.Marshal(ctx.AgentState); err == nil {
			if agentStateTable, err := p.jsonToLuaTable(string(agentStateJSON)); err == nil {
				p.vm.SetGlobal("agent_state", agentStateTable)
			}
		}
	}

	// Registry and Config are exposed as empty tables since they're complex Go objects
	// Plugins should use the execute function to interact with the system
	p.vm.SetGlobal("registry", p.vm.NewTable())
	p.vm.SetGlobal("config", p.vm.NewTable())

	initFn := p.vm.GetGlobal("initialize")
	if initFn != lua.LNil {
		if err := p.vm.CallByParam(lua.P{
			Fn:      initFn,
			NRet:    1,
			Protect: true,
		}); err != nil {
			return fmt.Errorf("initialize function failed: %w", err)
		}

		ret := p.vm.Get(-1)
		p.vm.Pop(1)

		if ret != lua.LNil && ret != lua.LFalse {
			if errStr, ok := ret.(lua.LString); ok {
				return fmt.Errorf("initialize failed: %s", string(errStr))
			}
		}
	}

	p.status = PluginStatusActive
	p.metadata.Status = PluginStatusActive
	logging.Info("Initialized Lua plugin: %s", p.name)
	return nil
}

func (p *LuaPlugin) Execute(cmd models.Command) models.CommandResult {
	executeFn := p.vm.GetGlobal("execute")
	if executeFn == lua.LNil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     "Plugin does not implement execute function",
			ResultType: models.ResultTypeError,
		}
	}

	cmdJSON, _ := json.Marshal(cmd)
	cmdTable, err := p.jsonToLuaTable(string(cmdJSON))
	if err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Failed to convert command: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	if err := p.vm.CallByParam(lua.P{
		Fn:      executeFn,
		NRet:    1,
		Protect: true,
	}, cmdTable); err != nil {
		return models.CommandResult{
			ID:         cmd.ID,
			Status:     models.CommandStatusFailed,
			Output:     fmt.Sprintf("Execute function error: %v", err),
			ResultType: models.ResultTypeError,
		}
	}

	ret := p.vm.Get(-1)
	p.vm.Pop(1)

	result := p.luaTableToCommandResult(ret, cmd.ID)
	return result
}

func (p *LuaPlugin) Cleanup() error {
	cleanupFn := p.vm.GetGlobal("cleanup")
	if cleanupFn != lua.LNil {
		if err := p.vm.CallByParam(lua.P{
			Fn:      cleanupFn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			return fmt.Errorf("cleanup function failed: %w", err)
		}
	}

	p.status = PluginStatusUnloaded
	p.metadata.Status = PluginStatusUnloaded
	p.context = nil
	return nil
}

func (p *LuaPlugin) GetStatus() PluginStatus {
	return p.status
}

func (p *LuaPlugin) jsonToLuaTable(jsonStr string) (*lua.LTable, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}
	return p.goValueToLuaValue(data).(*lua.LTable), nil
}

func (p *LuaPlugin) goValueToLuaValue(val interface{}) lua.LValue {
	switch v := val.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := p.vm.NewTable()
		for i, item := range v {
			table.RawSetInt(i+1, p.goValueToLuaValue(item))
		}
		return table
	case map[string]interface{}:
		table := p.vm.NewTable()
		for key, item := range v {
			table.RawSetString(key, p.goValueToLuaValue(item))
		}
		return table
	default:
		return lua.LNil
	}
}

func (p *LuaPlugin) luaTableToCommandResult(val lua.LValue, cmdID string) models.CommandResult {
	result := models.CommandResult{
		ID:         cmdID,
		Status:     models.CommandStatusCompleted,
		ResultType: models.ResultTypeText,
	}

	if val == lua.LNil {
		result.Status = models.CommandStatusFailed
		result.Output = "Plugin returned nil"
		return result
	}

	table, ok := val.(*lua.LTable)
	if !ok {
		result.Output = val.String()
		return result
	}

	if status := table.RawGetString("status"); status != lua.LNil {
		result.Status = models.CommandStatus(status.String())
	}
	if output := table.RawGetString("output"); output != lua.LNil {
		result.Output = output.String()
	}
	if resultType := table.RawGetString("result_type"); resultType != lua.LNil {
		result.ResultType = models.ResultType(resultType.String())
	}

	return result
}
