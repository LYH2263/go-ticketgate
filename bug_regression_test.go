package ticketgate

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestBug02_VerifyAfterCloseNoPanic(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = g.Verify(tok)
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed after Close, got %v", err)
	}
}
