package nonce

import (
	"sync"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

type Cache struct {
	mu  sync.Mutex
	lru *LRU
	clk clock.Clock
	ttl time.Duration
}

func NewCache(clk clock.Clock, cap int, ttl time.Duration) *Cache {
	if clk == nil {
		clk = clock.Real{}
	}
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &Cache{lru: NewLRU(cap), clk: clk, ttl: ttl}
}

func (c *Cache) TTL() time.Duration { return c.ttl }

func (c *Cache) Seen(nonce string) bool {
	if nonce == "" {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Sweep(c.clk.Now())
	exp, ok := c.lru.Get(nonce)
	if !ok {
		return false
	}
	return c.clk.Now().Before(exp)
}

func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lru.Len()
}
