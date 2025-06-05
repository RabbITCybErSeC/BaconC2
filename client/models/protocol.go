package models

type ITransportProtocol interface {
	Initialize() error
	Register(agent Agent) error
	Beacon() (Command, error)
	BeaconWithResultRequest() (Command, bool, error)
	SendResult(agentID string, result CommandResult) error
	Close() error
}
