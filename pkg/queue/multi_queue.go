package queue

import "sync"

// MultiQueue is a thread-safe interface for agent-based queues
type MultiQueue[T QueueItem] interface {
	Add(ID string, item T) error
	Get(ID string) (T, bool)
	List(ID string) ([]T, error)
}

// MemoryMultiQueue is a generic in-memory queue implementation for multiple agents
type MemoryMultiQueue[T QueueItem] struct {
	queues   map[string][]T
	queuesMu sync.RWMutex
}

// NewMemoryMultiQueue creates a new generic in-memory multi-queue
func NewMemoryMultiQueue[T QueueItem]() *MemoryMultiQueue[T] {
	return &MemoryMultiQueue[T]{
		queues: make(map[string][]T),
	}
}

// Add appends an item to the specified agent's queue
func (q *MemoryMultiQueue[T]) Add(ID string, item T) error {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	if _, exists := q.queues[ID]; !exists {
		q.queues[ID] = make([]T, 0)
	}

	q.queues[ID] = append(q.queues[ID], item)
	return nil
}

// Get retrieves and removes the first item from the specified agent's queue
func (q *MemoryMultiQueue[T]) Get(ID string) (T, bool) {
	q.queuesMu.Lock()
	defer q.queuesMu.Unlock()

	var zero T
	items, exists := q.queues[ID]
	if !exists || len(items) == 0 {
		return zero, false
	}

	item := items[0]
	q.queues[ID] = items[1:]
	return item, true
}

// List returns a copy of all items in the specified agent's queue
func (q *MemoryMultiQueue[T]) List(ID string) ([]T, error) {
	q.queuesMu.RLock()
	defer q.queuesMu.RUnlock()

	items, exists := q.queues[ID]
	if !exists {
		return []T{}, nil
	}

	itemsCopy := make([]T, len(items))
	copy(itemsCopy, items)
	return itemsCopy, nil
}
