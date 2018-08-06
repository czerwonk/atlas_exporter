package probe

import (
	"sync"
	"time"
)

// Cache caches probe lookup results
type Cache struct {
	cache map[int]*cacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

type cacheItem struct {
	expires time.Time
	value   *Probe
}

// NewCache creates a probe cache
func NewCache(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, cache: make(map[int]*cacheItem)}
}

// Get retrieves a probe from the cache (if exists, else returns false)
func (c *Cache) Get(id int) (*Probe, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if p, found := c.cache[id]; found && time.Now().Before(p.expires) {
		return p.value, true
	}

	return nil, false
}

// Add adds a probe to the cache
func (c *Cache) Add(id int, p *Probe) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[id] = &cacheItem{expires: time.Now().Add(c.ttl), value: p}
}

// CleanUp removes expired cache items
func (c *Cache) CleanUp() int {
	expired := make([]int, 0)

	for k, v := range c.cache {
		if v.expires.Before(time.Now()) {
			expired = append(expired, k)
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, id := range expired {
		delete(c.cache, id)
	}

	return len(expired)
}
