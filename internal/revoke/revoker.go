package revoke

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
)

type Mode int

const (
	ModeExact Mode = iota
	ModeBloom
	ModeBoth
)

type Revoker struct {
	mode  Mode
	exact *Exact
	bloom *Bloom
}

func New(clk clock.Clock, mode Mode) *Revoker {
	r := &Revoker{mode: mode, exact: NewExact(clk)}
	if mode == ModeBloom || mode == ModeBoth {
		r.bloom = NewBloom(8192, 5)
	}
	return r
}

func (r *Revoker) Revoke(jti string, until time.Time) error {
	if jti == "" {
		return fmt.Errorf("revoke: empty jti")
	}
	switch r.mode {
	case ModeExact:
		r.exact.Revoke(jti, until)
	case ModeBloom:
		r.bloom.Add([]byte(jti))
	case ModeBoth:
		r.exact.Revoke(jti, until)
		r.bloom.Add([]byte(jti))
	}
	return nil
}

func (r *Revoker) IsRevoked(jti string) bool {
	if jti == "" {
		return false
	}
	switch r.mode {
	case ModeExact:
		return r.exact.IsRevoked(jti)
	case ModeBloom:
		return r.bloom.MightContain([]byte(jti))
	case ModeBoth:
		if r.exact.IsRevoked(jti) {
			return true
		}
		return r.bloom.MightContain([]byte(jti))
	default:
		return r.exact.IsRevoked(jti)
	}
}

func (r *Revoker) Mode() Mode { return r.mode }
