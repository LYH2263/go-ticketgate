package ticketgate

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

func newGW(t *testing.T, clk *clock.Fake, opts ...Option) *Gateway {
	t.Helper()
	base := []Option{
		WithClock(clk),
		WithIssuer("tg"),
		WithAudience("api"),
		WithTTL(15 * time.Minute),
		WithGrace(10 * time.Second),
		WithSkew(2 * time.Second),
		WithNonceCache(64, time.Hour),
	}
	g, err := New(append(base, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = g.Close() })
	return g
}

func TestIssueVerifyRoundtrip(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, kid, err := g.Issue(Claims{Subject: "u1", Scope: []string{"read"}, Attrs: map[string]string{"k": "v"}})
	if err != nil {
		t.Fatal(err)
	}
	if kid == "" || len(tok) < 40 {
		t.Fatalf("empty token kid=%q len=%d", kid, len(tok))
	}
	c, err := g.Verify(tok)
	if err != nil {
		t.Fatal(err)
	}
	if c.Subject != "u1" || c.Issuer != "tg" || !c.HasAudience("api") || !c.HasScope("read") {
		t.Fatalf("claims %+v", c)
	}
	if v, ok := c.Attr("k"); !ok || v != "v" {
		t.Fatalf("attrs %v", c.Attrs)
	}
	if c.ID == "" || c.Nonce == "" {
		t.Fatal("missing jti/nonce")
	}
}

func TestExpiry(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk, WithTTL(10*time.Second), WithSkew(2*time.Second))
	tok, _, err := g.Issue(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	clk.Advance(12 * time.Second)
	if _, err := g.Verify(tok); err != nil {
		t.Fatalf("at exp+skew should still be valid, got %v", err)
	}
	clk.Advance(time.Second)
	if _, err := g.Verify(tok); !errors.Is(err, ErrExpired) {
		t.Fatalf("want ErrExpired, got %v", err)
	}
}

func TestRevoke(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	issued, err := g.IssueDetailed(Claims{Subject: "u", ID: "jti-rev-1"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.Revoke(issued.Claims.ID); err != nil {
		t.Fatal(err)
	}
	if !g.IsRevoked(issued.Claims.ID) {
		t.Fatal("expected revoked")
	}
	if _, err := g.Verify(issued.Token); !errors.Is(err, ErrRevoked) {
		t.Fatalf("want ErrRevoked, got %v", err)
	}
}

func TestKidGrace(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk, WithGrace(10*time.Second))
	tok, oldKid, err := g.Issue(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	newKid, err := g.RotateKey()
	if err != nil {
		t.Fatal(err)
	}
	if newKid == oldKid {
		t.Fatal("kid did not change")
	}
	if _, err := g.Verify(tok); err != nil {
		t.Fatalf("within grace: %v", err)
	}
	clk.Advance(10 * time.Second)
	if _, err := g.Verify(tok); !errors.Is(err, ErrKeyRetired) {
		t.Fatalf("want ErrKeyRetired after grace, got %v", err)
	}
}

func TestAudMismatch(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u", Audience: []string{"other"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Verify(tok); !errors.Is(err, ErrAudience) {
		t.Fatalf("want ErrAudience, got %v", err)
	}
}

func TestNonceReplay(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u", Nonce: "n-fixed"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Verify(tok); err != nil {
		t.Fatal(err)
	}
	if _, err := g.Verify(tok); !errors.Is(err, ErrNonceReplay) {
		t.Fatalf("want ErrNonceReplay, got %v", err)
	}
}

func TestRotateTicket(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	old, err := g.IssueDetailed(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	nw, err := g.RotateTicket(old.Token, Claims{})
	if err != nil {
		t.Fatal(err)
	}
	if nw.Claims.ID == old.Claims.ID {
		t.Fatal("jti should change")
	}
	if _, err := g.Verify(old.Token); !errors.Is(err, ErrRevoked) && !errors.Is(err, ErrNonceReplay) {
		t.Fatalf("old token should fail, got %v", err)
	}
	if _, err := g.Verify(nw.Token); err != nil {
		t.Fatal(err)
	}
}

func TestNotYetValid(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk, WithSkew(0))
	tok, _, err := g.Issue(Claims{Subject: "u", NotBefore: clk.Now().Add(30 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Verify(tok); !errors.Is(err, ErrNotYetValid) {
		t.Fatalf("want ErrNotYetValid, got %v", err)
	}
	clk.Advance(30 * time.Second)
	if _, err := g.Verify(tok); err != nil {
		t.Fatal(err)
	}
}

func TestTamperMAC(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	g := newGW(t, clk)
	tok, _, err := g.Issue(Claims{Subject: "u"})
	if err != nil {
		t.Fatal(err)
	}
	tok[len(tok)-1] ^= 0xff
	if _, err := g.Verify(tok); !errors.Is(err, ErrBadMAC) {
		t.Fatalf("want ErrBadMAC, got %v", err)
	}
}
