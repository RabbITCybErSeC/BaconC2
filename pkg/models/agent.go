package models

type Agent struct {
	ID           string `json:"id" gorm:"column:id;primaryKey"`
	Hostname     string `json:"hostname" gorm:"column:hostname"`
	IP           string `json:"ip" gorm:"column:ip"`
	OS           string `json:"os" gorm:"column:os"`
	Protocol     string `json:"protocol" gorm:"column:protocol"`
	ExtendedInfo string `json:"extended_info,omitempty" gorm:"type:text"`
}

type ExtendedAgentInfo struct {
	AgentID           string             `json:"agent_id" gorm:"column:agent_id;primaryKey"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces" gorm:"type:json"`
	Architecture      string             `json:"architecture" gorm:"column:architecture"`
	CPUInfo           string             `json:"cpu_info" gorm:"column:cpu_info"`
	MemoryTotal       uint64             `json:"memory_total" gorm:"column:memory_total"`
	MemoryFree        uint64             `json:"memory_free" gorm:"column:memory_free"`
	DiskTotal         uint64             `json:"disk_total" gorm:"column:disk_total"`
	DiskFree          uint64             `json:"disk_free" gorm:"column:disk_free"`
	Uptime            uint64             `json:"uptime" gorm:"column:uptime"`
	ProcessCount      int                `json:"process_count" gorm:"column:process_count"`
	Username          string             `json:"username" gorm:"column:username"`
	Domain            string             `json:"domain" gorm:"column:domain"`
	LastBootTime      string             `json:"last_boot_time" gorm:"column:last_boot_time"`
}
