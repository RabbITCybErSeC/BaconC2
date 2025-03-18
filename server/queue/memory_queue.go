package queue

import (
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/server/models"
)

type MemoryCommandQueue struct {
	queues   map[string][]models.Command
	queuesMu sync.RWMutex
}

func NewMemoryCommandQueue() *MemoryCommandQueue {
	return &MemoryCommandQueue{
		queues: make(map[string][]models.Command),
	}
}

func (q *MemoryCommandQueue) Add(agentID string, cmd models.Command) error {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	if _, exists := q.queues[agentID]; !exists {
		q.queues[agentID] = []models.Command{}
	}

	q.queues[agentID] = append(q.queues[agentID], cmd)
	return nil
}

func (q *MemoryCommandQueue) Get(agentID string) (models.Command, bool) {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	var emptyCmd models.Command

	commands, exists := q.queues[agentID]
	if !exists || len(commands) == 0 {
		return emptyCmd, false
	}

	cmd := commands[0]

	q.queues[agentID] = commands[1:]

	return cmd, true
}

func (q *MemoryCommandQueue) List(agentID string) ([]models.Command, error) {
	q.queuesMu.RLock()
	defer q.queuesMu.RUnlock()

	commands, exists := q.queues[agentID]
	if !exists {
		return []models.Command{}, nil
	}

	cmdCopy := make([]models.Command, len(commands))
	copy(cmdCopy, commands)

	return cmdCopy, nil
}
