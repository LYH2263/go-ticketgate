package trace

import (
	"sync"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

type Event struct {
	Time  time.Time
	Op    string
	KID   string
	JTI   string
	Stage string
	Err   string
}

type Log struct {
	mu   sync.Mutex
	clk  clock.Clock
	cap  int
	ring []Event
	next int
	full bool
}

func New(clk clock.Clock, cap int) *Log {
	if clk == nil {
		clk = clock.Real{}
	}
	if cap <= 0 {
		cap = 64
	}
	return &Log{clk: clk, cap: cap, ring: make([]Event, cap)}
}

func (l *Log) Record(op, kid, jti, stage string, err error) {
	if l == nil {
		return
	}
	e := Event{Time: l.clk.Now(), Op: op, KID: kid, JTI: jti, Stage: stage}
	if err != nil {
		e.Err = err.Error()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ring[l.next] = e
	l.next = (l.next + 1) % l.cap
	if l.next == 0 {
		l.full = true
	}
}

func (l *Log) Recent(n int) []Event {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	total := l.next
	if l.full {
		total = l.cap
	}
	if n <= 0 || n > total {
		n = total
	}
	out := make([]Event, 0, n)
	start := l.next - n
	if l.full {
		start = l.next - n
	}
	for i := 0; i < n; i++ {
		idx := start + i
		for idx < 0 {
			idx += l.cap
		}
		idx = idx % l.cap
		out = append(out, l.ring[idx])
	}
	return out
}
