package issue

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Draft struct {
	Issuer    string
	Subject   string
	Audience  []string
	TTL       time.Duration
	NotBefore time.Time
	JTI       string
	Nonce     string
	Scope     []string
	Attrs     map[string]string
	Bind      []byte
}

func Build(now time.Time, pol Policy, d Draft) (token.Payload, error) {
	iss := d.Issuer
	if iss == "" {
		iss = pol.DefaultIssuer
	}
	aud := cloneAUD(d.Audience)
	if len(aud) == 0 {
		aud = cloneAUD(pol.DefaultAUD)
	}
	if pol.RequireAUD && len(aud) == 0 {
		return token.Payload{}, fmt.Errorf("issue: aud required")
	}
	if pol.MaxAudiences > 0 && len(aud) > pol.MaxAudiences {
		return token.Payload{}, fmt.Errorf("issue: too many aud")
	}
	if len(aud) > 8 {
		return token.Payload{}, fmt.Errorf("issue: too many aud")
	}
	if pol.RequireSub && d.Subject == "" {
		return token.Payload{}, fmt.Errorf("issue: sub required")
	}
	ttl, err := pol.ResolveTTL(d.TTL)
	if err != nil {
		return token.Payload{}, err
	}
	jti, err := MustID(d.JTI, NewJTI, pol.RequireJTI, "jti")
	if err != nil {
		return token.Payload{}, err
	}
	nonce, err := MustID(d.Nonce, NewNonce, pol.RequireNonce, "nonce")
	if err != nil {
		return token.Payload{}, err
	}
	nbf := d.NotBefore
	if nbf.IsZero() {
		nbf = now
	}
	iat := now.UTC()
	exp := iat.Add(ttl)
	attrs := map[string]string{}
	for k, v := range d.Attrs {
		attrs[k] = v
	}
	p := token.Payload{
		Issuer:    iss,
		Subject:   d.Subject,
		Audience:  aud,
		IssuedAt:  iat,
		NotBefore: nbf.UTC(),
		ExpiresAt: exp,
		JTI:       jti,
		Nonce:     nonce,
		Scope:     append([]string(nil), d.Scope...),
		Attrs:     attrs,
		Bind:      append([]byte(nil), d.Bind...),
	}
	return p, nil
}
