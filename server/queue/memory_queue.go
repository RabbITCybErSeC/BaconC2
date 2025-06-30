package queue

import (
	"sync"

	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
)

type MemoryCommandQueue struct {
	queues   map[string][]local_models.AgentCommand
	queuesMu sync.RWMutex
}

func NewMemoryCommandQueue() *MemoryCommandQueue {
	return &MemoryCommandQueue{
		queues: make(map[string][]local_models.AgentCommand),
	}
}

func (q *MemoryCommandQueue) Add(agentID string, cmd local_models.AgentCommand) error {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	if _, exists := q.queues[agentID]; !exists {
		q.queues[agentID] = []local_models.AgentCommand{}
	}

	q.queues[agentID] = append(q.queues[agentID], cmd)
	return nil
}

func (q *MemoryCommandQueue) Get(agentID string) (local_models.AgentCommand, bool) {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	var emptyCmd local_models.AgentCommand

	commands, exists := q.queues[agentID]
	if !exists || len(commands) == 0 {
		return emptyCmd, false
	}

	cmd := commands[0]

	q.queues[agentID] = commands[1:]

	return cmd, true
}

func (q *MemoryCommandQueue) List(agentID string) ([]local_models.AgentCommand, error) {
	q.queuesMu.RLock()
	defer q.queuesMu.RUnlock()

	commands, exists := q.queues[agentID]
	if !exists {
		return []local_models.AgentCommand{}, nil
	}

	cmdCopy := make([]local_models.AgentCommand, len(commands))
	copy(cmdCopy, commands)

	return cmdCopy, nil
}
