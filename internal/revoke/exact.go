package revoke

import (
	"sync"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

type entry struct {
	Until time.Time
}

type Exact struct {
	mu   sync.Mutex
	clk  clock.Clock
	items map[string]entry
}

func NewExact(clk clock.Clock) *Exact {
	if clk == nil {
		clk = clock.Real{}
	}
	return &Exact{clk: clk, items: make(map[string]entry)}
}

func (e *Exact) Revoke(jti string, until time.Time) {
	if jti == "" {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.items[jti] = entry{Until: until}
}

func (e *Exact) IsRevoked(jti string) bool {
	if jti == "" {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.clk.Now()
	ent, ok := e.items[jti]
	if !ok {
		return false
	}
	if !ent.Until.IsZero() && !now.Before(ent.Until) {
		delete(e.items, jti)
		return false
	}
	return true
}

func (e *Exact) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.items)
}

func (e *Exact) Forget(jti string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.items, jti)
}
