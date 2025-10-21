package command_handler

import (
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

var commandRegistry *CommandHandlerRegistry

func init() {
	commandRegistry = NewCommandHandlerRegistry()
	// setup default built-in handlers
	commandRegistry.RegisterHandler(*NewGetRegisteredClientHandlers())
}

func GetGlobalCommandRegistry() *CommandHandlerRegistry {
	return commandRegistry
}

type CommandHandler struct {
	Name    string
	Handler func(cmd models.Command) models.CommandResult
}

type CommandHandlerRegistry struct {
	mu       sync.RWMutex
	handlers map[string]CommandHandler
}

func NewCommandHandlerRegistry() *CommandHandlerRegistry {
	return &CommandHandlerRegistry{
		handlers: make(map[string]CommandHandler),
	}
}

func (r *CommandHandlerRegistry) RegisterHandler(handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[handler.Name] = handler
}

func (r *CommandHandlerRegistry) GetHandler(name string) (CommandHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[name]
	return h, ok
}

func (r *CommandHandlerRegistry) GetAllRegisteredHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.handlers))
	for n := range r.handlers {
		names = append(names, n)
	}
	return names
}

func (r *CommandHandlerRegistry) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.handlers[name]; !exists {
		return false
	}
	delete(r.handlers, name)
	return true
}
