package nonce

import "time"

func (c *Cache) Occupy(nonce string, until time.Time) bool {
	if nonce == "" {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clk.Now()
	c.lru.Sweep(now)
	if exp, ok := c.lru.Get(nonce); ok && now.Before(exp) {
		return false
	}
	if until.IsZero() {
		until = now.Add(c.ttl)
	}
	c.lru.Put(nonce, until)
	return true
}

func (c *Cache) OccupyTTL(nonce string) bool {
	return c.Occupy(nonce, time.Time{})
}

func (c *Cache) Forget(nonce string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Delete(nonce)
	delete(c.pending, nonce)
}

func (c *Cache) Reserve(nonce string, until time.Time) bool {
	if nonce == "" {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clk.Now()
	if c.pending == nil {
		c.pending = map[string]time.Time{}
	}
	if exp, ok := c.pending[nonce]; ok && now.Before(exp) {
		return false
	}
	if until.IsZero() {
		until = now.Add(c.ttl)
	}
	c.pending[nonce] = until
	return true
}

func (c *Cache) Release(nonce string) {
	if nonce == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.pending, nonce)
}
