package command_handler

import (
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/formatter"
)

func NewGetRegisteredClientHandlers() *CommandHandler {
	return &CommandHandler{
		Name:    "get_registered_handlers",
		Handler: GetRegisteredClientHandlers,
	}
}

func GetRegisteredClientHandlers(cmd models.Command) models.CommandResult {
	registry := GetGlobalCommandRegistry()
	handlerNames := registry.GetAllRegisteredHandlers()

	outputData := formatter.ToJsonString(handlerNames)

	return models.CommandResult{
		ID:         cmd.ID,
		Status:     models.CommandStatusCompleted,
		Output:     outputData,
		ResultType: models.ResultTypeList,
	}
}
