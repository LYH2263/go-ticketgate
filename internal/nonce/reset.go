package nonce

func (c *Cache) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lru.cap
}

func (c *Cache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru = NewLRU(c.lru.cap)
}
