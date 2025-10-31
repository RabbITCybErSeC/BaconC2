package command_handler

import "github.com/RabbITCybErSeC/BaconC2/pkg/models"

type IAgentState interface {
	GetWorkingDirectory() string
	SetWorkingDirectory(path string) error
	GetEnv(key string) (string, bool)
	SetEnv(key, value string)
	GetAllEnv() map[string]string
}

type CommandContext struct {
	Command models.Command
	State   IAgentState
}

type StatefulCommandHandler struct {
	Name    string
	Handler func(ctx *CommandContext) models.CommandResult
}

func NewCommandContext(cmd models.Command, state IAgentState) *CommandContext {
	return &CommandContext{
		Command: cmd,
		State:   state,
	}
}
