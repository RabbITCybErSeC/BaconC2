package models

import "github.com/google/uuid"

type CommandStatus string
type CommandType string

const (
	CommandTypeInternal CommandType = "intern"
	CommandTypeShell    CommandType = "shell"
)

const (
	CommandStatusPending            CommandStatus = "cs_pndg"
	CommandStatusRunning            CommandStatus = "cs_rng"
	CommandStatusCompleted          CommandStatus = "cs_cmpltd"
	CommandStatusFailed             CommandStatus = "cs_fld"
	CommandStatusCancelled          CommandStatus = "cs_clld"
	CommandStatusTimeout            CommandStatus = "cs_tmt"
	CommandStatusAck                CommandStatus = "cs_ack"
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
	Type    CommandType `json:"type"`
	Command string      `json:"command"`
	Args    []string    `json:"args,omitempty"`
}

type Command struct {
	ID      string        `json:"id" gorm:"primaryKey"`
	Command string        `json:"command"`
	Args    []string      `json:"args,omitempty" gorm:"type:json"`
	Type    CommandType   `json:"type"`
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

func NewWebSocketMessage(msgType, data, id string) *WebSocketMessage {
	return &WebSocketMessage{
		Type: msgType,
		Data: data,
		ID:   id,
	}
}

func NewRawCommand(cmdType CommandType, command string, args ...string) *RawCommand {
	return &RawCommand{
		Type:    cmdType,
		Command: command,
		Args:    args,
	}
}

func NewCommand(command string, cmdType CommandType, args ...string) *Command {
	return &Command{
		ID:      uuid.New().String(),
		Command: command,
		Args:    args,
		Type:    cmdType,
		Status:  CommandStatusPending,
	}
}

func NewCommandResult(commandID string, status CommandStatus, output string) *CommandResult {
	return &CommandResult{
		ID:     commandID,
		Status: status,
		Output: output,
	}
}
