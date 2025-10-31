package filesystem

import (
	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

func PwdHandler(ctx *command_handler.CommandContext) models.CommandResult {
	currentDir := ctx.State.GetWorkingDirectory()

	return models.CommandResult{
		ID:     ctx.Command.ID,
		Status: models.CommandStatusCompleted,
		Output: currentDir,
	}
}

func NewPwdHandler() command_handler.StatefulCommandHandler {
	return command_handler.StatefulCommandHandler{
		Name:    "pwd",
		Handler: PwdHandler,
	}
}
