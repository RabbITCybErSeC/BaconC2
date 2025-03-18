package models

type ITransportProtocol interface {
	Initialize() error
	Register(agent Agent) error
	Beacon() (Command, error)
	SendResult(agentID string, result Command) error
	Close() error
}
