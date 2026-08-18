package issue

import (
	"fmt"
	"time"
)

type Policy struct {
	DefaultTTL    time.Duration
	MinTTL        time.Duration
	MaxTTL        time.Duration
	DefaultIssuer string
	RequireAUD    bool
	RequireSub    bool
	RequireNonce  bool
	RequireJTI    bool
	MaxAudiences  int
	DefaultAUD    []string
}

func DefaultPolicy() Policy {
	return Policy{
		DefaultTTL:    15 * time.Minute,
		MinTTL:        time.Second,
		MaxTTL:        24 * time.Hour,
		DefaultIssuer: "ticketgate",
		RequireAUD:    true,
		RequireSub:    false,
		RequireNonce:  true,
		RequireJTI:    true,
		MaxAudiences:  8,
	}
}

func (p Policy) ResolveTTL(explicit time.Duration) (time.Duration, error) {
	ttl := explicit
	if ttl <= 0 {
		ttl = p.DefaultTTL
	}
	if p.MinTTL > 0 && ttl < p.MinTTL {
		return 0, fmt.Errorf("issue: ttl below min")
	}
	if p.MaxTTL > 0 && ttl > p.MaxTTL {
		return 0, fmt.Errorf("issue: ttl above max")
	}
	if ttl <= 0 {
		return 0, fmt.Errorf("issue: ttl")
	}
	return ttl, nil
}

func cloneAUD(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
