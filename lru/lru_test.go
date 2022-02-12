package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("123123"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "123123" {
		t.Fatalf("cache hit key1=123123 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "val1", "val2", "val3"
	maxBytes := len(k1 + k2 + v1 + v2)
	lru := New(int64(maxBytes), callback)

	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	expect := []string{"key1", "key2"}

	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("call OnEvicte3d failed,expect: %s,but: %s", expect, keys)
	}
}
