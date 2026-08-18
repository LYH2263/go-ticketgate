package ticketgate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestBug04_VerifyHonorsCanceledContext(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, err := g.VerifyContext(ctx, tok)
	if err == nil {
		t.Fatalf("canceled verify must not return claims, got %+v", c)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
}
