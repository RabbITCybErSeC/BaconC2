package models

type CommandStatus string

const (
	CommandStatusPending            CommandStatus = "pending"
	CommandStatusRunning            CommandStatus = "running"
	CommandStatusCompleted          CommandStatus = "completed"
	CommandStatusFailed             CommandStatus = "failed"
	CommandStatusCancelled          CommandStatus = "cancelled"
	CommandStatusTimeout            CommandStatus = "timeout"
	CommandStatusAck                CommandStatus = "ack"
	CommandStatusSentToClient       CommandStatus = "c_sent"
	CommandStatusSentToServer       CommandStatus = "s_sent"
	CommandStatusReceivedFromClient CommandStatus = "c_received"
	CommandStatusReceivedFromServer CommandStatus = "s_received"
)

type WebSocketMessage struct {
	Type string `json:"type"` // "input", "output", "error", "control"
	Data string `json:"data"`
	ID   string `json:"id,omitempty"`
}

type RawCommand struct {
	Command string `json:"command"`
}

type Command struct {
	ID      string        `json:"id" gorm:"primaryKey"`
	Command string        `json:"command"`
	Status  CommandStatus `json:"status"`
}

type CommandResult struct {
	ID     string        `json:"id"`
	Status CommandStatus `json:"status"`
	Output string        `json:"output,omitempty" gorm:"type:text"`
}

type ICommandExecutor interface {
	Execute(cmd Command) CommandResult
	ProcessCommandQueue()
}
