package dist

import (
	"hash/crc64"
	"sync"
)

type shard[V any] struct {
	sync.RWMutex
	m map[string]V
}

type ShardedMap[V any] struct {
	shardNum uint64
	shards   []*shard[V]
}

func NewShardedMap[V any](shardNum uint64) *ShardedMap[V] {
	shards := make([]*shard[V], shardNum)
	for i := range shards {
		shards[i] = &shard[V]{
			m: make(map[string]V),
		}
	}
	return &ShardedMap[V]{
		shardNum: shardNum,
		shards:   shards,
	}
}

func (sm *ShardedMap[V]) Get(key string) (V, bool) {
	shard := sm.getShard(key)
	sm.shards[shard].RLock()
	defer sm.shards[shard].RUnlock()
	val, ok := sm.shards[shard].m[key]
	return val, ok
}

func (sm *ShardedMap[V]) Set(key string, val V) {
	shard := sm.getShard(key)
	sm.shards[shard].Lock()
	defer sm.shards[shard].Unlock()
	sm.shards[shard].m[key] = val
}

func (sm *ShardedMap[V]) Delete(key string) {
	shard := sm.getShard(key)
	sm.shards[shard].Lock()
	defer sm.shards[shard].Unlock()
	delete(sm.shards[shard].m, key)
}

func (sm *ShardedMap[V]) Keys() []string {
	keys := make([]string, 0, len(sm.shards))
	mut := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(sm.shards))

	for shard := range sm.shards {
		go func(shard int) {
			sm.shards[shard].RLock()
			for key := range sm.shards[shard].m {
				mut.Lock()
				keys = append(keys, key)
				mut.Unlock()
			}
			sm.shards[shard].RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return keys
}

func (sm *ShardedMap[V]) getShard(key string) uint64 {
	return crc64.Checksum([]byte(key), crc64.MakeTable(crc64.ECMA)) % uint64(sm.shardNum)
}
