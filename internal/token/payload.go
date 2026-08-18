package token

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/codec"
)

type Payload struct {
	Issuer    string
	Subject   string
	Audience  []string
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	JTI       string
	Nonce     string
	Scope     []string
	Attrs     map[string]string
	Bind      []byte
}

func unixSec(t time.Time) uint64 {
	if t.IsZero() {
		return 0
	}
	return uint64(t.Unix())
}

func fromUnix(sec uint64) time.Time {
	if sec == 0 {
		return time.Time{}
	}
	return time.Unix(int64(sec), 0).UTC()
}

func EncodePayload(p Payload) ([]byte, error) {
	if err := validatePayload(p); err != nil {
		return nil, err
	}
	m := codec.NewMap()
	if p.Issuer != "" {
		m.AddString(codec.FIss, p.Issuer)
	}
	if p.Subject != "" {
		m.AddString(codec.FSub, p.Subject)
	}
	for _, a := range p.Audience {
		m.AddString(codec.FAud, a)
	}
	if !p.IssuedAt.IsZero() {
		var b [8]byte
		codec.PutU64(b[:], unixSec(p.IssuedAt))
		m.Add(codec.FIat, b[:])
	}
	if !p.NotBefore.IsZero() {
		var b [8]byte
		codec.PutU64(b[:], unixSec(p.NotBefore))
		m.Add(codec.FNbf, b[:])
	}
	if !p.ExpiresAt.IsZero() {
		var b [8]byte
		codec.PutU64(b[:], unixSec(p.ExpiresAt))
		m.Add(codec.FExp, b[:])
	}
	if p.JTI != "" {
		m.AddString(codec.FJTI, p.JTI)
	}
	if p.Nonce != "" {
		m.AddString(codec.FNonce, p.Nonce)
	}
	for _, s := range p.Scope {
		m.AddString(codec.FScope, s)
	}
	if len(p.Attrs) > 0 {
		nb, err := codec.EncodeAttrs(p.Attrs)
		if err != nil {
			return nil, err
		}
		m.Add(codec.FAttrs, nb)
	}
	if len(p.Bind) > 0 {
		m.Add(codec.FBind, p.Bind)
	}
	return codec.EncodeMap(m, codec.MaxPayloadBytes)
}

func DecodePayload(b []byte) (Payload, error) {
	var p Payload
	m, err := codec.DecodeMap(b, false, true)
	if err != nil {
		return p, err
	}
	p.Issuer = m.FirstString(codec.FIss)
	p.Subject = m.FirstString(codec.FSub)
	p.Audience = m.AllString(codec.FAud)
	if v, ok := m.First(codec.FIat); ok {
		sec, err := codec.AsU64(v)
		if err != nil {
			return p, err
		}
		p.IssuedAt = fromUnix(sec)
	}
	if v, ok := m.First(codec.FNbf); ok {
		sec, err := codec.AsU64(v)
		if err != nil {
			return p, err
		}
		p.NotBefore = fromUnix(sec)
	}
	if v, ok := m.First(codec.FExp); ok {
		sec, err := codec.AsU64(v)
		if err != nil {
			return p, err
		}
		p.ExpiresAt = fromUnix(sec)
	}
	p.JTI = m.FirstString(codec.FJTI)
	p.Nonce = m.FirstString(codec.FNonce)
	p.Scope = m.AllString(codec.FScope)
	if v, ok := m.First(codec.FAttrs); ok {
		attrs, err := codec.DecodeAttrs(v)
		if err != nil {
			return p, err
		}
		p.Attrs = attrs
	}
	if v, ok := m.First(codec.FBind); ok {
		p.Bind = v
	}
	if err := validatePayload(p); err != nil {
		return p, err
	}
	return p, nil
}

func validatePayload(p Payload) error {
	if len(p.Issuer) > codec.MaxIssuerLen {
		return fmt.Errorf("token: iss too long")
	}
	if len(p.Subject) > codec.MaxSubjectLen {
		return fmt.Errorf("token: sub too long")
	}
	if len(p.Audience) > codec.MaxAudiences {
		return fmt.Errorf("token: too many aud")
	}
	for _, a := range p.Audience {
		if a == "" || len(a) > codec.MaxAudLen {
			return fmt.Errorf("token: aud")
		}
	}
	if len(p.JTI) > codec.MaxJTIBytes {
		return fmt.Errorf("token: jti")
	}
	if len(p.Nonce) > codec.MaxNonceBytes {
		return fmt.Errorf("token: nonce")
	}
	if len(p.Scope) > codec.MaxScopes {
		return fmt.Errorf("token: too many scopes")
	}
	for _, s := range p.Scope {
		if s == "" || len(s) > codec.MaxScopeLen {
			return fmt.Errorf("token: scope")
		}
	}
	if len(p.Bind) > codec.MaxBindBytes {
		return fmt.Errorf("token: bind")
	}
	if len(p.Attrs) > codec.MaxAttrs {
		return fmt.Errorf("token: attrs")
	}
	return nil
}
