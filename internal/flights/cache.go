package flights

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, val any, ttl time.Duration)
	Clear()
}

type InMemoryTTL struct {
	mu sync.RWMutex
	m  map[string]entry
}

type entry struct {
	v   any
	exp time.Time
}

func NewInMemoryTTL() *InMemoryTTL {
	return &InMemoryTTL{m: make(map[string]entry)}
}

func (c *InMemoryTTL) Get(k string) (any, bool) {
	c.mu.RLock()
	e, ok := c.m[k]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.exp) {
		return nil, false
	}
	return e.v, true
}

func (c *InMemoryTTL) Set(k string, v any, ttl time.Duration) {
	c.mu.Lock()
	c.m[k] = entry{v: v, exp: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *InMemoryTTL) Clear() {
	c.mu.Lock()
	c.m = make(map[string]entry)
	c.mu.Unlock()
}

// StartCleanup starts a goroutine that periodically removes expired entries
func (c *InMemoryTTL) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			c.mu.Lock()
			for k, e := range c.m {
				if now.After(e.exp) {
					delete(c.m, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
