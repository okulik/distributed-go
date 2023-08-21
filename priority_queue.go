package dist

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type PriorityQueueItem[K comparable] struct {
	Key       K
	Priority  int64
	Timestamp int64
	lru       int64
}

func NewPriorityQueueItem[K comparable](key K, priority int64, ts int64) *PriorityQueueItem[K] {
	return &PriorityQueueItem[K]{
		Key:       key,
		Priority:  priority,
		Timestamp: ts,
		lru:       ((priority & 0x00000000ffffffff) << 32) | (ts & 0xffffffff),
	}
}

func (pq *PriorityQueueItem[K]) String() string {
	return fmt.Sprintf("k: %v, p: %d, ts: %d", pq.Key, pq.Priority, pq.Timestamp)
}

var ErrNoFreeSlots = errors.New("no free slots in the queue")
var ErrEmptyQueue = errors.New("no items in the queue")
var ErrNoSuchKey = errors.New("no item with such key in the queue")

type PriorityQueue[K comparable] struct {
	capacity int
	ind      map[K]int
	data     []*PriorityQueueItem[K]
	mu       sync.RWMutex
}

func NewPriorityQueue[K comparable](capacity int) *PriorityQueue[K] {
	return &PriorityQueue[K]{
		capacity: capacity,
		ind:      make(map[K]int, capacity),
		data:     make([]*PriorityQueueItem[K], 0, capacity),
	}
}

func (pq *PriorityQueue[K]) Push(key K, priority int64, ts int64) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.data) == pq.capacity {
		return ErrNoFreeSlots
	}

	pq.data = append(pq.data, NewPriorityQueueItem[K](key, priority, ts))
	pq.ind[key] = len(pq.data) - 1
	pq.heapify_up(len(pq.data) - 1)

	return nil
}

func (pq *PriorityQueue[K]) Pop() (*PriorityQueueItem[K], error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.data) == 0 {
		return nil, ErrEmptyQueue
	}

	if len(pq.data) > 1 {
		pq.data[0], pq.data[len(pq.data)-1] = pq.data[len(pq.data)-1], pq.data[0]
	}

	item := pq.data[len(pq.data)-1]
	pq.data = pq.data[:len(pq.data)-1]
	delete(pq.ind, item.Key)
	pq.heapify_down(0)

	return item, nil
}

func (pq *PriorityQueue[K]) Peek() (*PriorityQueueItem[K], error) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.data) == 0 {
		return nil, ErrEmptyQueue
	}
	return pq.data[0], nil
}

func (pq *PriorityQueue[K]) RemoveAt(key K) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.data) == 0 {
		return ErrEmptyQueue
	}

	i, ok := pq.ind[key]
	if !ok {
		return ErrNoSuchKey
	}
	if i != len(pq.data)-1 {
		pq.data[i] = pq.data[len(pq.data)-1]
		pq.ind[pq.data[i].Key] = i
		pq.data = pq.data[:len(pq.data)-1]
		delete(pq.ind, key)
		pq.heapify_down(i)
		return nil
	}

	delete(pq.ind, key)
	pq.data = pq.data[:len(pq.data)-1]

	return nil
}

func (pq *PriorityQueue[K]) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.data)
}

func (pq *PriorityQueue[K]) String() string {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	res := make([]string, 0, len(pq.data))
	for _, item := range pq.data {
		res = append(res, item.String())
	}
	return strings.Join(res, ", ")
}

func (pq *PriorityQueue[K]) heapify_up(i int) {
	if len(pq.data) <= 1 || i == 0 {
		return
	}

	parent := pq.parent(i)
	if pq.data[parent].lru <= pq.data[i].lru {
		return
	}

	pq.ind[pq.data[parent].Key], pq.ind[pq.data[i].Key] = i, parent
	pq.data[parent], pq.data[i] = pq.data[i], pq.data[parent]
	pq.heapify_up(parent)
}

func (pq *PriorityQueue[K]) heapify_down(i int) {
	if len(pq.data) <= 1 {
		return
	}

	left := pq.left(i)
	right := pq.right(i)

	var min int
	if left < len(pq.data) && pq.data[left].lru <= pq.data[i].lru {
		min = left
	} else {
		min = i
	}

	if right < len(pq.data) && pq.data[right].lru <= pq.data[min].lru {
		min = right
	}

	if min != i {
		pq.ind[pq.data[i].Key], pq.ind[pq.data[min].Key] = min, i
		pq.data[i], pq.data[min] = pq.data[min], pq.data[i]
		pq.heapify_down(min)
	}
}

func (pq *PriorityQueue[K]) left(i int) int {
	return i*2 + 1
}

func (pq *PriorityQueue[K]) right(i int) int {
	return i*2 + 2
}

func (pq *PriorityQueue[K]) parent(i int) int {
	return (i - 1) / 2
}
