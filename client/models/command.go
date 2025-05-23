package models

type Command struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Status  string `json:"status"`
	Output  string `json:"output,omitempty"`
}

type ICommandExecutor interface {
	Execute(cmd Command) Command
	ProcessCommandQueue()
}
