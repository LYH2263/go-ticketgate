package ticketgate

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/issue"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/revoke"
	"github.com/LYH2263/go-ticketgate/internal/skew"
	"github.com/LYH2263/go-ticketgate/internal/token"
	"github.com/LYH2263/go-ticketgate/internal/trace"
	"github.com/LYH2263/go-ticketgate/internal/verify"
)

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	switch {
	case contains(s, "bad magic"):
		return ErrBadMagic
	case contains(s, "unsupported version"):
		return ErrUnsupportedVersion
	case contains(s, "truncated"):
		return ErrTruncated
	case contains(s, "too large"):
		return ErrTooLarge
	case contains(s, "field order"):
		return ErrFieldOrder
	case contains(s, "unknown field"):
		return ErrUnknownField
	case contains(s, "duplicate"):
		return ErrDuplicateField
	case contains(s, "bad mac"):
		return ErrBadMAC
	case contains(s, "unknown kid"):
		return ErrUnknownKID
	case contains(s, "key retired"):
		return ErrKeyRetired
	case contains(s, "expired"):
		return ErrExpired
	case contains(s, "not yet valid"):
		return ErrNotYetValid
	case contains(s, "audience mismatch"):
		return ErrAudience
	case contains(s, "issuer mismatch"):
		return ErrIssuer
	case contains(s, "subject mismatch"):
		return ErrSubject
	case contains(s, "revoked"):
		return ErrRevoked
	case contains(s, "nonce replay"):
		return ErrNonceReplay
	case contains(s, "no current key"):
		return ErrNoCurrentKey
	case contains(s, "ttl"):
		return ErrTTL
	case contains(s, "required") || contains(s, "aud required") || contains(s, "sub required"):
		return ErrRequired
	case contains(s, "scope"):
		return ErrScope
	case contains(s, "algorithm") || contains(s, "alg"):
		return ErrAlg
	case contains(s, "empty jti"):
		return ErrEmptyJTI
	case contains(s, "bind"):
		return ErrBind
	default:
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	n, m := len(s), len(sub)
	if m == 0 {
		return 0
	}
	for i := 0; i+m <= n; i++ {
		if s[i:i+m] == sub {
			return i
		}
	}
	return -1
}

func claimsFrom(p token.Payload) Claims {
	attrs := map[string]string{}
	for k, v := range p.Attrs {
		attrs[k] = v
	}
	return Claims{
		Issuer:    p.Issuer,
		Subject:   p.Subject,
		Audience:  append([]string(nil), p.Audience...),
		IssuedAt:  p.IssuedAt,
		NotBefore: p.NotBefore,
		ExpiresAt: p.ExpiresAt,
		ID:        p.JTI,
		Nonce:     p.Nonce,
		Scope:     append([]string(nil), p.Scope...),
		Attrs:     attrs,
		Bind:      append([]byte(nil), p.Bind...),
	}
}

func draftFrom(c Claims) issue.Draft {
	ttl := time.Duration(0)
	if !c.IssuedAt.IsZero() && !c.ExpiresAt.IsZero() && c.ExpiresAt.After(c.IssuedAt) {
		ttl = c.ExpiresAt.Sub(c.IssuedAt)
	}
	return issue.Draft{
		Issuer:    c.Issuer,
		Subject:   c.Subject,
		Audience:  append([]string(nil), c.Audience...),
		TTL:       ttl,
		NotBefore: c.NotBefore,
		JTI:       c.ID,
		Nonce:     c.Nonce,
		Scope:     append([]string(nil), c.Scope...),
		Attrs:     c.Attrs,
		Bind:      c.Bind,
	}
}

func newParts(cfg config) (*keyring.Ring, *issue.Issuer, *verify.Verifier, *revoke.Revoker, *nonce.Cache, *trace.Log, error) {
	ring, err := keyring.New(cfg.clk, cfg.grace, cfg.alg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	rev := revoke.New(cfg.clk, cfg.revokeMode)
	nc := nonce.NewCache(cfg.clk, cfg.nonceCap, cfg.nonceTTL)
	iss := issue.New(ring, cfg.clk, cfg.policy())
	iss.SetNonceCache(nc)
	ver := verify.New(verify.Deps{
		Ring:    ring,
		Revoker: rev,
		Nonces:  nc,
		Skew:    skew.New(cfg.skew),
		Clock:   cfg.clk,
		Expect: verify.Expect{
			Audience: cfg.constraint(),
			Scopes:   append([]string(nil), cfg.scopes...),
		},
	})
	tr := trace.New(cfg.clk, cfg.traceCap)
	return ring, iss, ver, rev, nc, tr, nil
}
