package queue

import (
	"fmt"

	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
)

type CommandQueue interface {
	Add(agentID string, cmd local_models.ServerAgentModel) error
	Get(agentID string) (local_models.ServerAgentModel, bool)
	List(agentID string) ([]local_models.ServerAgentModel, error)
}

var ErrQueueEmpty = fmt.Errorf("command queue is empty")
