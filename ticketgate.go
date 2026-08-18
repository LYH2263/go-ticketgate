package ticketgate

import (
	"sync"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/issue"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/revoke"
	"github.com/LYH2263/go-ticketgate/internal/token"
	"github.com/LYH2263/go-ticketgate/internal/trace"
	"github.com/LYH2263/go-ticketgate/internal/verify"
)

// Gateway 会话票据网关：签发、验证、密钥轮换、票据轮换、吊销。
type Gateway struct {
	mu      sync.Mutex
	closed  bool
	ring    *keyring.Ring
	issuer  *issue.Issuer
	ver     *verify.Verifier
	revoker *revoke.Revoker
	nonces  *nonce.Cache
	trace   *trace.Log
}

func New(opts ...Option) (*Gateway, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if !cfg.alg.Valid() {
		cfg.alg = token.AlgHS256
	}
	ring, iss, ver, rev, nc, tr, err := newParts(cfg)
	if err != nil {
		return nil, err
	}
	return &Gateway{
		ring:    ring,
		issuer:  iss,
		ver:     ver,
		revoker: rev,
		nonces:  nc,
		trace:   tr,
	}, nil
}

func (g *Gateway) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.closed = true
	return nil
}

func (g *Gateway) Issue(claims Claims) ([]byte, string, error) {
	issued, err := g.IssueDetailed(claims)
	if err != nil {
		return nil, "", err
	}
	return issued.Token, issued.KID, nil
}

func (g *Gateway) IssueDetailed(claims Claims) (Issued, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return Issued{}, ErrClosed
	}
	raw, kid, p, err := g.issuer.Issue(draftFrom(claims))
	if err != nil {
		g.trace.Record("issue", kid, claims.ID, "", err)
		return Issued{}, mapErr(err)
	}
	out := Issued{Token: raw, KID: kid, Claims: claimsFrom(p)}
	g.trace.Record("issue", kid, out.Claims.ID, "ok", nil)
	return out, nil
}

func (g *Gateway) Verify(raw []byte) (Claims, error) {
	res, err := g.VerifyDetailed(raw)
	if err != nil {
		return Claims{}, err
	}
	return res.Claims, nil
}

func (g *Gateway) VerifyDetailed(raw []byte) (VerifyResult, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return VerifyResult{}, ErrClosed
	}
	o, err := g.ver.Verify(raw)
	if err != nil {
		g.trace.Record("verify", o.Header.KID, o.Payload.JTI, string(o.Stage), err)
		return VerifyResult{}, mapErr(err)
	}
	g.trace.Record("verify", o.Header.KID, o.Payload.JTI, "ok", nil)
	return VerifyResult{
		Claims: claimsFrom(o.Payload),
		KID:    o.Header.KID,
		Alg:    byte(o.Header.Alg),
	}, nil
}

func (g *Gateway) Revoke(jti string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return ErrClosed
	}
	if jti == "" {
		return ErrEmptyJTI
	}
	err := g.revoker.Revoke(jti, time.Time{})
	g.trace.Record("revoke", "", jti, "", err)
	return err
}

func (g *Gateway) IsRevoked(jti string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.revoker.IsRevoked(jti)
}

func (g *Gateway) RotateKey() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return "", ErrClosed
	}
	kid, err := g.ring.Rotate()
	g.trace.Record("rotate-key", kid, "", "", err)
	if err != nil {
		return "", mapErr(err)
	}
	return kid, nil
}

func (g *Gateway) GracePeriod() time.Duration {
	return g.ring.GracePeriod()
}

func (g *Gateway) Keys() []KeyInfo {
	list := g.ring.List()
	out := make([]KeyInfo, 0, len(list))
	for _, k := range list {
		out = append(out, KeyInfo{
			KID:      k.Material.KID,
			Alg:      byte(k.Material.Alg),
			Status:   k.Status.String(),
			Created:  k.Created,
			RetireAt: k.RetireAt,
		})
	}
	return out
}

func (g *Gateway) Recent(n int) []string {
	ev := g.trace.Recent(n)
	out := make([]string, 0, len(ev))
	for _, e := range ev {
		line := e.Op + " " + e.KID + " " + e.JTI + " " + e.Stage
		if e.Err != "" {
			line += " " + e.Err
		}
		out = append(out, line)
	}
	return out
}
