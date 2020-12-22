package pcache

import (
	"pcache/pcache/lru"
	"pcache/pcache/view"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes uint64
}

func (c *cache) Add(key string, val view.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, val)
}

func (c *cache) Get(key string) (value view.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if val, ok := c.lru.Get(key); ok {
		return val.(view.ByteView), true
	}
	return
}
