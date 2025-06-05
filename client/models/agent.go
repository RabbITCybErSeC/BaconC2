package models

type Agent struct {
	ID           string             `json:"id"`
	Hostname     string             `json:"hostname"`
	IP           string             `json:"ip"`
	OS           string             `json:"os"`
	Protocol     string             `json:"protocol"`
	ExtendedInfo *ExtendedAgentInfo `json:"extended_info,omitempty"`
}

type ExtendedAgentInfo struct {
	NetworkInterfaces []NetworkInterface `json:"network_interfaces"`
	Architecture      string             `json:"architecture"`
	CPUInfo           string             `json:"cpu_info"`
	MemoryTotal       uint64             `json:"memory_total"`
	MemoryFree        uint64             `json:"memory_free"`
	DiskTotal         uint64             `json:"disk_total"`
	DiskFree          uint64             `json:"disk_free"`
	Uptime            uint64             `json:"uptime"`
	ProcessCount      int                `json:"process_count"`
	Username          string             `json:"username"`
	Domain            string             `json:"domain"`
	LastBootTime      string             `json:"last_boot_time"`
}
