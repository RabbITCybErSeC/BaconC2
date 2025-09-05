package queue

import (
	"fmt"
	"sync"
)

var ErrQueueEmpty = fmt.Errorf("queue is empty")
var ErrItemNotFound = fmt.Errorf("item not found")
var ErrIndexOutOfRange = fmt.Errorf("index out of range")

type QueueItem interface {
	any
}

type GenericQueue[T QueueItem] interface {
	Add(item T) error
	Get() (T, bool)
	List() ([]T, error)
	RemoveFirst() (T, error)
	RemoveAt(index int) (T, error)
	RemoveItem(item T) error
	Size() int
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

func (q *MemoryQueue[T]) RemoveFirst() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	var zero T
	if len(q.items) == 0 {
		return zero, ErrQueueEmpty
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *MemoryQueue[T]) RemoveAt(index int) (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	var zero T
	if index < 0 || index >= len(q.items) {
		return zero, ErrIndexOutOfRange
	}
	item := q.items[index]
	q.items = append(q.items[:index], q.items[index+1:]...)
	return item, nil
}

func (q *MemoryQueue[T]) RemoveItem(item T) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, queueItem := range q.items {
		// Note: This uses == comparison, only works for comparable types
		// For more complex comparison, we might need to pass a comparison function in the future
		if any(queueItem) == any(item) {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}

func (q *MemoryQueue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}
