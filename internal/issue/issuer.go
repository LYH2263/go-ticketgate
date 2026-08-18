package issue

import (
	"fmt"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Issuer struct {
	ring   *keyring.Ring
	clk    clock.Clock
	pol    Policy
	nonces *nonce.Cache
}

func New(ring *keyring.Ring, clk clock.Clock, pol Policy) *Issuer {
	if clk == nil {
		clk = clock.Real{}
	}
	if pol.DefaultTTL == 0 {
		pol = DefaultPolicy()
	}
	return &Issuer{ring: ring, clk: clk, pol: pol}
}

func (i *Issuer) SetNonceCache(nc *nonce.Cache) { i.nonces = nc }

func (i *Issuer) Policy() Policy { return i.pol }

func (i *Issuer) Issue(d Draft) ([]byte, string, token.Payload, error) {
	if i.ring == nil {
		return nil, "", token.Payload{}, fmt.Errorf("issue: no keyring")
	}
	p, err := Build(i.clk.Now(), i.pol, d)
	if err != nil {
		return nil, "", token.Payload{}, err
	}
	if i.nonces != nil && p.Nonce != "" {
		if !i.nonces.Reserve(p.Nonce, p.ExpiresAt) {
			return nil, "", token.Payload{}, fmt.Errorf("issue: nonce replay")
		}
	}
	cur, err := i.ring.Current()
	if err != nil {
		if i.nonces != nil && p.Nonce != "" {
			i.nonces.Release(p.Nonce)
		}
		return nil, "", token.Payload{}, err
	}
	raw, err := token.IssueBytes(cur.Material.KID, cur.Material.Alg, cur.Material.Secret, p)
	if err != nil {
		if i.nonces != nil && p.Nonce != "" {
			i.nonces.Release(p.Nonce)
		}
		return nil, "", token.Payload{}, err
	}
	if i.nonces != nil && p.Nonce != "" {
		i.nonces.Release(p.Nonce)
	}
	return raw, cur.Material.KID, p, nil
}
