package ticketgate

import (
	"time"

	"github.com/LYH2263/go-ticketgate/internal/audience"
	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/issue"
	"github.com/LYH2263/go-ticketgate/internal/revoke"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

type config struct {
	clk        clock.Clock
	grace      time.Duration
	skew       time.Duration
	ttl        time.Duration
	minTTL     time.Duration
	maxTTL     time.Duration
	issuer     string
	aud        []string
	subject    string
	alg        token.Alg
	nonceCap   int
	nonceTTL   time.Duration
	revokeMode revoke.Mode
	requireSub bool
	scopes     []string
	traceCap   int
	useGlob    bool
}

func defaultConfig() config {
	return config{
		clk:        clock.Real{},
		grace:      time.Hour,
		skew:       0,
		ttl:        15 * time.Minute,
		minTTL:     time.Second,
		maxTTL:     24 * time.Hour,
		issuer:     "ticketgate",
		alg:        token.AlgHS256,
		nonceCap:   4096,
		nonceTTL:   24 * time.Hour,
		revokeMode: revoke.ModeExact,
		traceCap:   128,
	}
}

type Option func(*config)

func WithClock(c clock.Clock) Option {
	return func(cfg *config) {
		if c != nil {
			cfg.clk = c
		}
	}
}

func WithGrace(d time.Duration) Option {
	return func(cfg *config) { cfg.grace = d }
}

func WithSkew(d time.Duration) Option {
	return func(cfg *config) { cfg.skew = d }
}

func WithTTL(d time.Duration) Option {
	return func(cfg *config) { cfg.ttl = d }
}

func WithTTLBounds(min, max time.Duration) Option {
	return func(cfg *config) {
		cfg.minTTL = min
		cfg.maxTTL = max
	}
}

func WithIssuer(iss string) Option {
	return func(cfg *config) { cfg.issuer = iss }
}

func WithAudience(aud ...string) Option {
	return func(cfg *config) { cfg.aud = append([]string(nil), aud...) }
}

func WithSubject(sub string) Option {
	return func(cfg *config) { cfg.subject = sub }
}

func WithAlg(alg byte) Option {
	return func(cfg *config) { cfg.alg = token.Alg(alg) }
}

func WithNonceCache(cap int, ttl time.Duration) Option {
	return func(cfg *config) {
		cfg.nonceCap = cap
		cfg.nonceTTL = ttl
	}
}

func WithRevokeBloom() Option {
	return func(cfg *config) { cfg.revokeMode = revoke.ModeBloom }
}

func WithRevokeBoth() Option {
	return func(cfg *config) { cfg.revokeMode = revoke.ModeBoth }
}

func WithRequireSubject() Option {
	return func(cfg *config) { cfg.requireSub = true }
}

func WithScopes(s ...string) Option {
	return func(cfg *config) { cfg.scopes = append([]string(nil), s...) }
}

func WithGlobAudience() Option {
	return func(cfg *config) { cfg.useGlob = true }
}

func WithTraceCap(n int) Option {
	return func(cfg *config) { cfg.traceCap = n }
}

func (c config) policy() issue.Policy {
	return issue.Policy{
		DefaultTTL:    c.ttl,
		MinTTL:        c.minTTL,
		MaxTTL:        c.maxTTL,
		DefaultIssuer: c.issuer,
		RequireAUD:    len(c.aud) > 0,
		RequireSub:    c.requireSub,
		RequireNonce:  true,
		RequireJTI:    true,
		MaxAudiences:  8,
		DefaultAUD:    append([]string(nil), c.aud...),
	}
}

func (c config) constraint() audience.Constraint {
	return audience.Constraint{
		Issuer:   c.issuer,
		Audience: append([]string(nil), c.aud...),
		Subject:  c.subject,
		UseGlob:  c.useGlob,
	}
}
