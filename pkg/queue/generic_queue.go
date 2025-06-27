package queue

import (
	"fmt"
	"sync"
)

var ErrQueueEmpty = fmt.Errorf("queue is empty")

type QueueItem interface {
	any
}

type GenericQueue[T QueueItem] interface {
	Add(item T) error
	Get() (T, bool)
	List() ([]T, error)
}
type MemoryQueue[T QueueItem] struct {
	items []T
	mu    sync.RWMutex
}

func NewMemoryQueue[T QueueItem]() *MemoryQueue[T] {
	return &MemoryQueue[T]{
		items: make([]T, 0),
	}
}

func (q *MemoryQueue[T]) Add(item T) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	return nil
}

func (q *MemoryQueue[T]) Get() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var zero T
	if len(q.items) == 0 {
		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *MemoryQueue[T]) List() ([]T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	itemsCopy := make([]T, len(q.items))
	copy(itemsCopy, q.items)
	return itemsCopy, nil
}
