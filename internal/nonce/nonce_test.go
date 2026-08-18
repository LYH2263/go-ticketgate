package nonce

import (
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestOccupyReplay(t *testing.T) {
	clk := clock.NewFake(time.Unix(10, 0))
	c := NewCache(clk, 8, time.Minute)
	if !c.Occupy("n1", clk.Now().Add(time.Minute)) {
		t.Fatal("first occupy")
	}
	if c.Occupy("n1", clk.Now().Add(time.Minute)) {
		t.Fatal("replay should fail")
	}
	clk.Advance(2 * time.Minute)
	if !c.Occupy("n1", clk.Now().Add(time.Minute)) {
		t.Fatal("after ttl should succeed")
	}
}

func TestLRUEvict(t *testing.T) {
	clk := clock.NewFake(time.Unix(10, 0))
	c := NewCache(clk, 2, time.Hour)
	_ = c.Occupy("a", clk.Now().Add(time.Hour))
	_ = c.Occupy("b", clk.Now().Add(time.Hour))
	_ = c.Occupy("c", clk.Now().Add(time.Hour))
	if c.Len() != 2 {
		t.Fatalf("len=%d", c.Len())
	}
}
