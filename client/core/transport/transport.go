package transport

import "github.com/RabbITCybErSeC/BaconC2/pkg/models"

type ITransportProtocol interface {
	Initialize(agent models.Agent) error
	Start() error
	Stop() error
	SendResults(cmd models.Command) models.CommandResult
}
