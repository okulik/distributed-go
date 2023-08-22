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
	ecb      EvictedCB[K, V]
}

type EvictedCB[K comparable, V any] func(key K, item *CacheItem[V])

var ErrLRUInternal = fmt.Errorf("lru internal error")

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
		ecb:      nil,
	}
}

func NewLRUCacheWithEvict[K comparable, V any](capacity int, onEvicted EvictedCB[K, V]) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		data:     make(map[K]*CacheItem[V]),
		lru:      NewPriorityQueue[K](capacity),
		exp:      NewPriorityQueue[K](capacity),
		ecb:      onEvicted,
	}
}

func (lc *LRUCache[K, V]) Get(key K) (*CacheItem[V], bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	ci, ok := lc.data[key]
	if ok {
		_ = lc.lru.RemoveAt(key)
		_ = lc.lru.Push(key, ci.Priority, time.Now().Unix())
	}
	return ci, ok
}

func (lc *LRUCache[K, V]) Add(key K, value V, priority int64, expiry int64) (evicted bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, ok := lc.data[key]
	if ok {
		_ = lc.lru.RemoveAt(key)
		_ = lc.lru.Push(key, priority, time.Now().Unix())
		_ = lc.exp.RemoveAt(key)
		_ = lc.exp.Push(key, expiry, 0)
		lc.data[key] = &CacheItem[V]{
			Value:    value,
			Priority: priority,
			Expiry:   expiry,
		}
		return false
	}

	// cleanup strategy remove all expired items
	if len(lc.data) == lc.capacity {
		item, err := lc.exp.Peek()
		for err == nil && item.Priority < time.Now().Unix() {
			expiredItem, _ := lc.exp.Pop()
			lc.purge(expiredItem.Key)
			evicted = true
			_ = lc.lru.RemoveAt(expiredItem.Key)
			item, err = lc.exp.Peek()
		}
	}

	// cleanup strategy remove lowest priority item
	if len(lc.data) == lc.capacity {
		lowestPrioItem, _ := lc.lru.Pop()
		lc.purge(lowestPrioItem.Key)
		evicted = true
		_ = lc.exp.RemoveAt(lowestPrioItem.Key)
	}

	lc.data[key] = &CacheItem[V]{
		Value:    value,
		Priority: priority,
		Expiry:   expiry,
	}
	_ = lc.lru.Push(key, priority, time.Now().Unix())
	_ = lc.exp.Push(key, expiry, 0)

	return evicted
}

func (lc *LRUCache[K, V]) Remove(key K) (present bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if _, ok := lc.data[key]; !ok {
		return false
	}

	delete(lc.data, key)
	_ = lc.lru.RemoveAt(key)
	_ = lc.exp.RemoveAt(key)
	return true
}

func (lc *LRUCache[K, V]) Contains(key K) bool {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	_, ok := lc.data[key]
	return ok
}

func (lc *LRUCache[K, V]) Len() int {
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

func (lc *LRUCache[K, V]) purge(key K) {
	evictedItem := lc.data[key]
	delete(lc.data, key)
	if lc.ecb != nil {
		lc.ecb(key, evictedItem)
	}
}
