package clock

import (
	"sync"
	"time"
)

// Clock 可注入时钟，便于过期与宽限期边界测试。
type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time { return time.Now() }

// Fake 可前进的测试时钟。
type Fake struct {
	mu sync.Mutex
	t  time.Time
}

func NewFake(t time.Time) *Fake {
	if t.IsZero() {
		t = time.Unix(1_700_000_000, 0).UTC()
	}
	return &Fake{t: t.UTC()}
}

func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

func (f *Fake) Set(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = t.UTC()
}

func (f *Fake) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = f.t.Add(d)
}
