package dist

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type LRUCache[K comparable, V any] struct {
	capacity int
	data     map[K]*CacheItem[V]
	lru      *PriorityQueue[K]
	exp      *PriorityQueue[K]
	mu       sync.RWMutex
}

type CacheItem[V any] struct {
	Value    V
	Priority int64
	Expiry   int64
}

func (ci CacheItem[V]) String() string {
	return fmt.Sprintf("value: %v, priority: %d, expiry: %d", ci.Value, ci.Priority, ci.Expiry)
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		data:     make(map[K]*CacheItem[V]),
		lru:      NewPriorityQueue[K](capacity),
		exp:      NewPriorityQueue[K](capacity),
	}
}

func (lc *LRUCache[K, V]) Get(key K) (*CacheItem[V], bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	ci, ok := lc.data[key]
	if ok {
		lc.lru.RemoveAt(key)
		lc.lru.Push(key, ci.Priority, time.Now().Unix())
	}
	return ci, ok
}

func (lc *LRUCache[K, V]) Set(key K, value V, priority int64, expiry int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, ok := lc.data[key]
	if ok {
		lc.lru.RemoveAt(key)
		lc.lru.Push(key, priority, time.Now().Unix())
		lc.exp.RemoveAt(key)
		lc.exp.Push(key, expiry, 0)
		lc.data[key] = &CacheItem[V]{
			Value:    value,
			Priority: priority,
			Expiry:   expiry,
		}
		return
	}

	// cleanup strategy remove all expired items
	if len(lc.data) == lc.capacity {
		item, err := lc.exp.Peek()
		for err == nil && item.Priority < time.Now().Unix() {
			expItem, _ := lc.exp.Pop()
			remKey := expItem.Key
			delete(lc.data, remKey)
			lc.lru.RemoveAt(remKey)
			item, err = lc.exp.Peek()
		}
	}

	// cleanup strategy remove lowest priority item
	if len(lc.data) == lc.capacity {
		item, _ := lc.lru.Pop()
		remKey := item.Key
		lc.exp.RemoveAt(remKey)
		delete(lc.data, remKey)
	}

	lc.data[key] = &CacheItem[V]{
		Value:    value,
		Priority: priority,
		Expiry:   expiry,
	}
	lc.lru.Push(key, priority, time.Now().Unix())
	lc.exp.Push(key, expiry, 0)
}

func (lc *LRUCache[K, V]) Size() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	return len(lc.data)
}

func (lc *LRUCache[K, V]) String() string {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	res := make([]string, 0, len(lc.data))
	for key, item := range lc.data {
		res = append(res, fmt.Sprintf("key: '%v', %s", key, item.String()))
	}
	return strings.Join(res, "\n")
}
