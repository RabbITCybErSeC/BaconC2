package models

import (
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type StreamingConfig struct {
	ShellType ShellType
	Term      string
}

func NewStreamingConfig(shellType ShellType, term string) *StreamingConfig {
	return &StreamingConfig{
		ShellType: shellType,
		Term:      term,
	}
}

type Capabilities struct {
	SupportsStreaming bool
	SupportsPush      bool
	SupportsBeacon    bool
}

type ITransportProtocol interface {
	Initialize(agent models.Agent) error
	RunProtocol() error
	// Beacon() (models.Command, error)
	// BeaconWithResultRequest() (models.Command, bool, error)
	// SendResults() error
	Close() error
}

type IStreamingTransport interface {
	StartStreamingSession(sessionType string, config *StreamingConfig, resultChan chan<- models.CommandResult) error
	CloseSession(sessionID string) error
}
