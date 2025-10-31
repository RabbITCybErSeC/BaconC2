package filesystem

import (
	"fmt"
	"path/filepath"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

// CdHandler changes the current working directory
func CdHandler(ctx *command_handler.CommandContext) models.CommandResult {
	var targetPath string

	if len(ctx.Command.Args) == 0 {
		targetPath = getHomeDirectory()
	} else {
		targetPath = ctx.Command.Args[0]
	}

	currentDir := ctx.State.GetWorkingDirectory()

	var absPath string
	if filepath.IsAbs(targetPath) {
		absPath = filepath.Clean(targetPath)
	} else {
		absPath = filepath.Clean(filepath.Join(currentDir, targetPath))
	}

	if err := ctx.State.SetWorkingDirectory(absPath); err != nil {
		return models.CommandResult{
			ID:     ctx.Command.ID,
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("Failed to change directory to '%s': %v", targetPath, err),
		}
	}

	newDir := ctx.State.GetWorkingDirectory()
	return models.CommandResult{
		ID:     ctx.Command.ID,
		Status: models.CommandStatusCompleted,
		Output: fmt.Sprintf("Changed directory to: %s", newDir),
	}
}

func NewCdHandler() command_handler.StatefulCommandHandler {
	return command_handler.StatefulCommandHandler{
		Name:    "cd",
		Handler: CdHandler,
	}
}
