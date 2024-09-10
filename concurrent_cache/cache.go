package cache

import (
	"distributed_cache/lru_cache" // Ensure the path is correct
	"sync"
)

// Define a struct instead of an interface
type Cache struct {
	mu         sync.Mutex
	lru        *lru_cache.Cache
	cacheBytes int64
}

func (c *Cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// lazy initialization
	if c.lru == nil {
		c.lru = lru_cache.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

