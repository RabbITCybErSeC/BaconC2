package models

type Command struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Status  string `json:"status"`
	Output  any    `json:"output,omitempty"`
}

type CommandResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output any    `json:"output"`
}

type ICommandExecutor interface {
	Execute(cmd Command) CommandResult
	ProcessCommandQueue()
}
