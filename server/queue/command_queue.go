package queue

import (
	"fmt"

	"github.com/RabbITCybErSeC/BaconC2/server/models"
)

type CommandQueue interface {
	Add(agentID string, cmd models.Command) error
	Get(agentID string) (models.Command, bool)
	List(agentID string) ([]models.Command, error)
}

var ErrQueueEmpty = fmt.Errorf("command queue is empty")
