package models

type WebSocketMessage struct {
	Type string `json:"type"` // "input", "output", "error", "control"
	Data string `json:"data"`
	ID   string `json:"id,omitempty"`
}

type Command struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Command string `json:"command"`
	Status  string `json:"status"`
}

type CommandResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output any    `json:"output,omitempty" gorm:"-"` // Excluded from database
}

type ICommandExecutor interface {
	Execute(cmd Command) CommandResult
	ProcessCommandQueue()
}
