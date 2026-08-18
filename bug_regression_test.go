package ticketgate

import (
	"strings"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/codec"
)

func TestBug03_IssueReleasesNonceOnFailure(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	stuck := "n-stuck-partial"
	longIss := strings.Repeat("x", codec.MaxIssuerLen+8)
	_, _, err := g.Issue(Claims{Subject: "u", Issuer: longIss, Nonce: stuck})
	if err == nil {
		t.Fatal("expected late issue failure")
	}
	tok, _, err := g.Issue(Claims{Subject: "u", Nonce: stuck})
	if err != nil {
		t.Fatalf("nonce must be released after failed issue: %v", err)
	}
	if _, err := g.Verify(tok); err != nil {
		t.Fatal(err)
	}
}
