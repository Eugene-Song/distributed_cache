package cache

import (
	"distributed_cache/lru_cache"
	"sync"
)

// Getter is the interface for getting the data
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc is a function type that satisfies the Getter interface
type GetterFunc func(key string) ([]byte, error)

// Get implements the Getter interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache *lru_cache.Cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup creates a new Group instance
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: lru_cache.New(cacheBytes, nil),
	}
	groups[name] = group
	return group
}

// GetGroup returns the named group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// Get is the main function for getting the data
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, nil
	}
	if v, ok := g.mainCache.Get(key); ok {
		return v.(ByteView), nil
	}
	return g.load(key)
}

// load is the function for loading the data
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// getLocally is the function for getting the data locally
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache is the function for adding the data to the cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value)
}