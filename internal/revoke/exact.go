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

// Get returns the current entry for jti, if any. Used to capture prior state
// so a revoke whose persist later fails can be rolled back exactly.
func (e *Exact) Get(jti string) (entry, bool) {
	if jti == "" {
		return entry{}, false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	ent, ok := e.items[jti]
	return ent, ok
}

// Restore puts the entry for jti back to ent, or removes it when had is false.
// It is the inverse of Revoke for the purpose of rolling back a failed persist.
func (e *Exact) Restore(jti string, ent entry, had bool) {
	if jti == "" {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if had {
		e.items[jti] = ent
	} else {
		delete(e.items, jti)
	}
}

func (e *Exact) List() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]string, 0, len(e.items))
	now := e.clk.Now()
	for jti, ent := range e.items {
		if !ent.Until.IsZero() && !now.Before(ent.Until) {
			continue
		}
		out = append(out, jti)
	}
	return out
}

func (e *Exact) Replace(jtis []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.items = make(map[string]entry, len(jtis))
	for _, jti := range jtis {
		if jti != "" {
			e.items[jti] = entry{}
		}
	}
}
