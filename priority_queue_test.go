package dist_test

import (
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

func TestPriorityQueuePush(t *testing.T) {
	t.Parallel()
	pq := PriorityQueueFixture(t)

	item, _ := pq.Peek()
	if item.Key != "b" {
		t.Errorf("Expected 'b', got '%v'", item.Key)
	}
}

func TestPriorityQueuePushFull(t *testing.T) {
	t.Parallel()
	pq := PriorityQueueFixture(t)

	err := pq.Push("d", 3000, time.Now().Unix())
	if err != dist.ErrNoFreeSlots {
		t.Errorf("Expected error dist.ErrNoFreeSlots")
	}
}

func TestPriorityQueuePop(t *testing.T) {
	t.Parallel()
	pq := PriorityQueueFixture(t)

	pqi, err := pq.Pop()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if pqi.Key != "b" {
		t.Errorf("Expected 'b', got '%v'", pqi.Key)
	}
}

func TestPriorityQueuePopEmpty(t *testing.T) {
	t.Parallel()
	pq := dist.NewPriorityQueue[string](3)

	_, err := pq.Pop()
	if err != dist.ErrEmptyQueue {
		t.Errorf("Expected error dist.ErrEmptyQueue, got nil")
	}
}

func TestPriorityQueueRemoveAt(t *testing.T) {
	t.Parallel()
	pq := PriorityQueueFixture(t)

	_ = pq.RemoveAt("b")
	item, _ := pq.Peek()
	if item.Key != "a" {
		t.Errorf("Expected 'a', got '%v'", item.Key)
	}
}

func TestPriorityQueueRemoveAtMissing(t *testing.T) {
	t.Parallel()
	pq := PriorityQueueFixture(t)

	err := pq.RemoveAt("d")
	if err != dist.ErrNoSuchKey {
		t.Errorf("Expected ErrNoSuchKey, got %v", err)
	}
}

func PriorityQueueFixture(t *testing.T) *dist.PriorityQueue[string] {
	pq := dist.NewPriorityQueue[string](3)
	_ = pq.Push("a", 2000, time.Now().Unix())
	_ = pq.Push("b", 1000, time.Now().Unix())
	_ = pq.Push("c", 3000, time.Now().Unix())
	return pq
}
