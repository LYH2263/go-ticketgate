package ticketgate

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/audience"
	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/issue"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/revoke"
	"github.com/LYH2263/go-ticketgate/internal/skew"
	"github.com/LYH2263/go-ticketgate/internal/token"
	"github.com/LYH2263/go-ticketgate/internal/verify"
)

type loadFailPersist struct {
	jtis    []string
	loadErr error
}

func (p *loadFailPersist) Save(jtis []string) error {
	p.jtis = append([]string(nil), jtis...)
	return nil
}

func (p *loadFailPersist) Load() ([]string, error) {
	if p.loadErr != nil {
		return nil, p.loadErr
	}
	return append([]string(nil), p.jtis...), nil
}

func TestBug10_VerifyFailClosedOnRevokeLoadError(t *testing.T) {
	clk := clock.NewFake(time.Unix(1_700_000_000, 0))
	ring, err := keyring.New(clk, time.Hour, token.AlgHS256)
	if err != nil {
		t.Fatal(err)
	}
	ps := &loadFailPersist{}
	rev1 := revoke.New(clk, revoke.ModeExact)
	rev1.SetPersist(ps)
	iss := issue.New(ring, clk, issue.DefaultPolicy().WithDefaultAUD("api"))
	raw, _, pld, err := iss.Issue(issue.Draft{Subject: "u", Audience: []string{"api"}, JTI: "jti-reload-1"})
	if err != nil {
		t.Fatal(err)
	}
	if err := rev1.Revoke(pld.JTI, time.Time{}); err != nil {
		t.Fatal(err)
	}
	ps.loadErr = errors.New("reload io")
	rev2 := revoke.New(clk, revoke.ModeExact)
	rev2.SetPersist(ps)
	ver := verify.New(verify.Deps{
		Ring:    ring,
		Revoker: rev2,
		Nonces:  nonce.NewCache(clk, 32, time.Hour),
		Skew:    skew.New(0),
		Clock:   clk,
		Expect:  verify.Expect{Audience: audience.Constraint{Issuer: "ticketgate", Audience: []string{"api"}}},
	})
	_, err = ver.Verify(raw)
	if err == nil {
		t.Fatal("verify must not accept a revoked ticket when reload fails")
	}
}
