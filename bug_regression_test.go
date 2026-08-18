package ticketgate

import (
	"bytes"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

func TestBug09_DecodeScratchNotAliased(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok1, _, err := g.Issue(Claims{Subject: "u1", Bind: []byte("bind-one-xxxx")})
	if err != nil {
		t.Fatal(err)
	}
	tok2, _, err := g.Issue(Claims{Subject: "u2", Bind: []byte("bind-two-yyyy")})
	if err != nil {
		t.Fatal(err)
	}
	env1, err := token.Split(tok1)
	if err != nil {
		t.Fatal(err)
	}
	env2, err := token.Split(tok2)
	if err != nil {
		t.Fatal(err)
	}
	p1, err := token.DecodePayload(env1.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p1.Bind, []byte("bind-one-xxxx")) {
		t.Fatalf("first bind %q", p1.Bind)
	}
	p2, err := token.DecodePayload(env2.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p2.Bind, []byte("bind-two-yyyy")) {
		t.Fatalf("second bind %q", p2.Bind)
	}
	if !bytes.Equal(p1.Bind, []byte("bind-one-xxxx")) {
		t.Fatalf("first payload bind aliased by second decode: %q", p1.Bind)
	}
}
