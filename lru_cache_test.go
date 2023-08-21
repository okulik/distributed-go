package dist_test

import (
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

func TestLRUCacheGet(t *testing.T) {
	t.Parallel()
	lc := LRUCacheFixture(t)

	v, ok := lc.Get("b")
	if !ok {
		t.Errorf("Item is missing")
	}
	if v.Value != 2 {
		t.Errorf("Expected 2, got %d", v)
	}
}

func TestLRUCacheSet(t *testing.T) {
	t.Parallel()
	lc := LRUCacheFixture(t)

	if lc.Size() != 3 {
		t.Errorf("Expected 3, got %d", lc.Size())
	}
}

func TestLRUCacheReplace(t *testing.T) {
	t.Parallel()
	lc := LRUCacheFixture(t)

	lc.Set("a", 4, 2000, time.Now().Unix()+30000)
	item, ok := lc.Get("a")
	if !ok {
		t.Errorf("Item is missing")
	}
	if item.Value != 4 {
		t.Errorf("Expected 4, got %d", item.Value)
	}
}

func TestLRUCacheSetFullRemoveLowestPrio(t *testing.T) {
	t.Parallel()
	lc := LRUCacheFixture(t)

	lc.Set("d", 4, 4000, time.Now().Unix()+30000)
	_, ok := lc.Get("b")
	if ok {
		t.Errorf("Expected item 'b' to be removed")
	}
}

func TestLRUCacheSetFullRemoveExpired(t *testing.T) {
	t.Parallel()
	lc := LRUCacheFixture(t)
	lc.Set("a", 1, 2000, time.Now().Unix()-10000)
	lc.Set("d", 4, 4000, time.Now().Unix()+30000)
	_, ok := lc.Get("a")
	if ok {
		t.Errorf("Expected item 'a' to be removed")
	}
	_, ok = lc.Get("d")
	if !ok {
		t.Errorf("Expected item 'd' to be present")
	}
}

func LRUCacheFixture(t *testing.T) *dist.LRUCache[string, int] {
	lc := dist.NewLRUCache[string, int](3)
	lc.Set("a", 1, 2000, time.Now().Unix()+30000)
	lc.Set("b", 2, 1000, time.Now().Unix()+30000)
	lc.Set("c", 3, 3000, time.Now().Unix()+30000)
	return lc
}
