package commands

import (
	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type CommandHandler func(cmd models.Command) models.CommandResult

type CommandHandlerRegistry struct {
	handlers map[string]CommandHandler
}

func NewCommandHandlerRegistry() *CommandHandlerRegistry {
	return &CommandHandlerRegistry{
		handlers: make(map[string]CommandHandler),
	}
}

func (r *CommandHandlerRegistry) RegisterHandler(command string, handler CommandHandler) {
	r.handlers[command] = handler
}

func (r *CommandHandlerRegistry) GetHandler(command string) (CommandHandler, bool) {
	handler, exists := r.handlers[command]
	return handler, exists
}
