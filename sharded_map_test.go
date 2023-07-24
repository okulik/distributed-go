package dist_test

import (
	"reflect"
	"testing"

	dist "github.com/okulik/distributed-go"
)

func TestShardedMap(t *testing.T) {
	smap := dist.NewShardedMap[int](3)

	smap.Set("create", 1)
	smap.Set("read", 2)
	smap.Set("update", 3)
	smap.Set("delete", 4)
	smap.Set("list", 5)
	smap.Set("get", 6)

	v, ok := smap.Get("create")
	if !ok {
		t.Fatal("missing key create")
	}
	if v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}

	v, ok = smap.Get("read")
	if !ok {
		t.Fatal("missing key read")
	}
	if v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	v, ok = smap.Get("update")
	if !ok {
		t.Fatal("missing key update")
	}
	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	v, ok = smap.Get("delete")
	if !ok {
		t.Fatal("missing key delete")
	}
	if v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}

	v, ok = smap.Get("list")
	if !ok {
		t.Fatal("missing key list")
	}
	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	v, ok = smap.Get("get")
	if !ok {
		t.Fatal("missing key get")
	}
	if v != 6 {
		t.Fatalf("expected 6, got %d", v)
	}

	keys := smap.Keys()
	if len(keys) != 6 {
		t.Fatalf("expected 6 keys, got %d", len(keys))
	}
	if reflect.DeepEqual(keys, []string{"create", "read", "update", "delete", "list", "get"}) {
		t.Fatalf("expected [create read update delete list get], got %v", keys)
	}
}
