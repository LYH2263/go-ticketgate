package ticketgate

import (
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/audience"
	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/issue"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/skew"
	"github.com/LYH2263/go-ticketgate/internal/token"
	"github.com/LYH2263/go-ticketgate/internal/verify"
)

func TestBug06_CurrentKeySecretNotAliased(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	ring, err := keyring.New(clk, time.Hour, token.AlgHS256)
	if err != nil {
		t.Fatal(err)
	}
	nc := nonce.NewCache(clk, 32, time.Hour)
	iss := issue.New(ring, clk, issue.DefaultPolicy().WithDefaultAUD("api"))
	iss.SetNonceCache(nc)
	raw, kid, _, err := iss.Issue(issue.Draft{Subject: "u", Audience: []string{"api"}})
	if err != nil {
		t.Fatal(err)
	}
	k, err := ring.Get(kid)
	if err != nil {
		t.Fatal(err)
	}
	if len(k.Material.Secret) == 0 {
		t.Fatal("empty secret")
	}
	k.Material.Secret[0] ^= 0xff
	ver := verify.New(verify.Deps{
		Ring:   ring,
		Nonces: nc,
		Skew:   skew.New(0),
		Clock:  clk,
		Expect: verify.Expect{Audience: audience.Constraint{Issuer: "ticketgate", Audience: []string{"api"}}},
	})
	if _, err := ver.Verify(raw); err != nil {
		t.Fatalf("mutating Get() secret must not change ring material: %v", err)
	}
}
