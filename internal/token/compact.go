package token

func CompactLen(header, payload, mac int) int {
	return 9 + header + payload + mac
}

func ValidMACLen(n int) bool {
	return n == 32 || n == 48 || n == 64
}

func ClonePayload(p Payload) Payload {
	attrs := map[string]string{}
	for k, v := range p.Attrs {
		attrs[k] = v
	}
	return Payload{
		Issuer:    p.Issuer,
		Subject:   p.Subject,
		Audience:  append([]string(nil), p.Audience...),
		IssuedAt:  p.IssuedAt,
		NotBefore: p.NotBefore,
		ExpiresAt: p.ExpiresAt,
		JTI:       p.JTI,
		Nonce:     p.Nonce,
		Scope:     append([]string(nil), p.Scope...),
		Attrs:     attrs,
		Bind:      append([]byte(nil), p.Bind...),
	}
}

func PayloadEqual(a, b Payload) bool {
	if a.Issuer != b.Issuer || a.Subject != b.Subject || a.JTI != b.JTI || a.Nonce != b.Nonce {
		return false
	}
	if unixSec(a.IssuedAt) != unixSec(b.IssuedAt) {
		return false
	}
	if unixSec(a.NotBefore) != unixSec(b.NotBefore) {
		return false
	}
	if unixSec(a.ExpiresAt) != unixSec(b.ExpiresAt) {
		return false
	}
	if len(a.Audience) != len(b.Audience) {
		return false
	}
	for i := range a.Audience {
		if a.Audience[i] != b.Audience[i] {
			return false
		}
	}
	if len(a.Scope) != len(b.Scope) {
		return false
	}
	for i := range a.Scope {
		if a.Scope[i] != b.Scope[i] {
			return false
		}
	}
	if len(a.Bind) != len(b.Bind) {
		return false
	}
	for i := range a.Bind {
		if a.Bind[i] != b.Bind[i] {
			return false
		}
	}
	if len(a.Attrs) != len(b.Attrs) {
		return false
	}
	for k, v := range a.Attrs {
		if b.Attrs[k] != v {
			return false
		}
	}
	return true
}
