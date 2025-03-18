package config

import "time"

type AgentConfig struct {
	ServerURL      string
	BeaconInterval time.Duration
	AgentID        string
	Protocol       string
	UDPServerHost  string
	UDPServerPort  int
}
