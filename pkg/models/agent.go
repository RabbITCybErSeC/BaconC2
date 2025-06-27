package models

type Agent struct {
	ID           string             `json:"id" gorm:"column:agent_id;primaryKey"`
	Hostname     string             `json:"hostname" gorm:"column:hostname"`
	IP           string             `json:"ip" gorm:"column:ip"`
	OS           string             `json:"os" gorm:"column:os"`
	Protocol     string             `json:"protocol" gorm:"column:protocol"`
	ExtendedInfo *ExtendedAgentInfo `json:"extended_info,omitempty" gorm:"column:extended_info;type:jsonb"`
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
