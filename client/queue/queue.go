package queue

import (
	"fmt"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type CommandQueue interface {
	Add(cmd models.Command) error
	Get() (models.Command, bool)
	List() ([]models.Command, error)
}

var ErrQueueEmpty = fmt.Errorf("command queue is empty")

type MemoryCommandQueue struct {
	commands []models.Command
	mu       sync.RWMutex
}

func NewMemoryCommandQueue() *MemoryCommandQueue {
	return &MemoryCommandQueue{
		commands: make([]models.Command, 0),
	}
}

func (q *MemoryCommandQueue) Add(cmd models.Command) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.commands = append(q.commands, cmd)
	return nil
}

func (q *MemoryCommandQueue) Get() (models.Command, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var emptyCmd models.Command

	if len(q.commands) == 0 {
		return emptyCmd, false
	}

	cmd := q.commands[0]
	q.commands = q.commands[1:]

	return cmd, true
}

func (q *MemoryCommandQueue) List() ([]models.Command, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	cmdCopy := make([]models.Command, len(q.commands))
	copy(cmdCopy, q.commands)

	return cmdCopy, nil
}
