package ticketgate

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

func TestBug05_DecodeWrapsCorruptSentinel(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	_, err := g.Verify([]byte("not-a-ticket"))
	if err == nil {
		t.Fatal("expected corrupt token error")
	}
	if !errors.Is(err, ErrCorrupt) {
		t.Fatalf("want errors.Is(ErrCorrupt), got %v", err)
	}
	_, err = token.Split([]byte("XXXX"))
	if !errors.Is(err, token.ErrCorrupt) {
		t.Fatalf("split want token.ErrCorrupt, got %v", err)
	}
}
