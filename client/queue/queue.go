package queue

import (
	"fmt"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type ICommandQueue interface {
	Add(cmd models.Command) error
	Get() (models.Command, bool)
	List() ([]models.Command, error)
}

type IResultQueue interface {
	Add(result models.CommandResult) error
	Get() (models.CommandResult, bool)
	List() ([]models.CommandResult, error)
}

var ErrQueueEmpty = fmt.Errorf("command queue is empty")

// MemoryCommandQueue is an in-memory implementation of CommandQueue
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

type MemoryResultQueue struct {
	results []models.CommandResult
	mu      sync.RWMutex
}

func NewMemoryResultQueue() *MemoryResultQueue {
	return &MemoryResultQueue{
		results: make([]models.CommandResult, 0),
	}
}

func (q *MemoryResultQueue) Add(result models.CommandResult) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.results = append(q.results, result)
	return nil
}

func (q *MemoryResultQueue) Get() (models.CommandResult, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var emptyResult models.CommandResult

	if len(q.results) == 0 {
		return emptyResult, false
	}

	result := q.results[0]
	q.results = q.results[1:]

	return result, true
}

func (q *MemoryResultQueue) List() ([]models.CommandResult, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	resultCopy := make([]models.CommandResult, len(q.results))
	copy(resultCopy, q.results)

	return resultCopy, nil
}
