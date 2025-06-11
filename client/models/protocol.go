package models

type ITransportProtocol interface {
	Initialize() error
	Register(agent Agent) error
	Beacon() (Command, error)
	BeaconWithResultRequest() (Command, bool, error)
	SendResult(agentID string, result CommandResult) error
	Close() error
}

type IStreamingTransport interface {
	StartStreamingSession(sessionType string, config map[string]interface{}, resultChan chan<- CommandResult) error
	CloseSession(sessionID string) error
}
