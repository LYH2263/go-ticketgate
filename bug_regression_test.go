package ticketgate

import (
	"bytes"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func TestBug01_ReturnedClaimsAliasTokenBytes(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u", Bind: []byte("bind-value-ok")})
	if err != nil {
		t.Fatal(err)
	}
	orig := append([]byte(nil), tok...)
	c, err := g.Verify(tok)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Bind) == 0 {
		t.Fatal("expected bind")
	}
	c.Bind[0] ^= 0xff
	if !bytes.Equal(tok, orig) {
		t.Fatal("mutating returned claims must not change the token bytes")
	}
}
