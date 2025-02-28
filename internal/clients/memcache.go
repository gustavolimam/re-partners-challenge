package clients

import (
	"math"
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.Mutex
}

// NewCache creates a new instance of Cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache with an expiration time in seconds
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).Unix()
	} else {
		expiration = math.MaxInt64 // Never expires
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get returns an item from the cache, if it is still valid
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.Expiration != math.MaxInt64 && item.Expiration < time.Now().Unix() {
		// Remove the expired item from the cache
		delete(c.items, key)
		return nil, false
	}
	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// CleanUp removes expired items from the cache
func (c *Cache) CleanUp() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().Unix()

	for key, item := range c.items {
		if item.Expiration != math.MaxInt64 && item.Expiration < now {
			delete(c.items, key)
		}
	}
}
