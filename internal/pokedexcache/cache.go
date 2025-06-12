package pokedexcache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt	time.Time
	val			[]byte
}

type Cache struct {
	Map		map[string]cacheEntry
	Mutex	sync.Mutex
	Interval	time.Duration
}

func NewCache(interval time.Duration) Cache {
	var mu sync.Mutex
	
	cache := Cache{
		Map: make(map[string]cacheEntry),
		Mutex: mu,
		Interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Map[key] = cacheEntry{
		createdAt:	time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	entry, exists := c.Map[key]
	if exists {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	ticker := time.NewTicker(c.Interval)
	for {
		<-ticker.C
		for key, val := range c.Map {
			if time.Since(val.createdAt) > c.Interval {
				delete(c.Map, key)
			}
		}
	}
}
