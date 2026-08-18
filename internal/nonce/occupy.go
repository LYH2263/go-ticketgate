package nonce

import "time"

// Occupy 在缓存窗内首次占位成功返回 true；已存在且未过期返回 false（重放）。
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
}
