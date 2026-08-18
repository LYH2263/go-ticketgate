package ticketgate

import (
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestBug07_RotateIssueThenRevoke(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	old, err := g.IssueDetailed(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	tooMany := []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8"}
	_, err = g.RotateTicket(old.Token, Claims{Audience: tooMany})
	if err == nil {
		t.Fatal("expected new issue to fail")
	}
	if g.IsRevoked(old.Claims.ID) {
		t.Fatal("must issue new ticket before revoking old")
	}
}
