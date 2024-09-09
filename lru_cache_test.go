package lru_cache

import (
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

// test cache basic add and get
func TestGet(t *testing.T) {
	lruCache := New(int64(0), nil)
	lruCache.Add("key1", String("1234"))
	// for test get a key in the cache
	if v, ok := lruCache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	// for test a key not in the cache
	if _, ok := lruCache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// test cache remove oldest element
func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + k2 + v1 + v2) // len is the number of bytes
	lruCache := New(int64(cap), nil)
	lruCache.Add(k1, String(v1))
	lruCache.Add(k2, String(v2))
	lruCache.Add(k3, String(v3))

	// shoule remove the oldest k1
	if _, ok := lruCache.Get("k1"); ok || lruCache.Len() != 2 {
		t.Fatalf("RemoveOldest k1 failed")
	}
}
