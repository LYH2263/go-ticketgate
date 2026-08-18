package ticketgate

import "time"

func (c Claims) Clone() Claims {
	attrs := map[string]string{}
	for k, v := range c.Attrs {
		attrs[k] = v
	}
	return Claims{
		Issuer:    c.Issuer,
		Subject:   c.Subject,
		Audience:  append([]string(nil), c.Audience...),
		IssuedAt:  c.IssuedAt,
		NotBefore: c.NotBefore,
		ExpiresAt: c.ExpiresAt,
		ID:        c.ID,
		Nonce:     c.Nonce,
		Scope:     append([]string(nil), c.Scope...),
		Attrs:     attrs,
		Bind:      append([]byte(nil), c.Bind...),
	}
}

func (c Claims) TTL() time.Duration {
	if c.IssuedAt.IsZero() || c.ExpiresAt.IsZero() {
		return 0
	}
	d := c.ExpiresAt.Sub(c.IssuedAt)
	if d < 0 {
		return 0
	}
	return d
}

func (c Claims) HasAudience(aud string) bool {
	for _, a := range c.Audience {
		if a == aud {
			return true
		}
	}
	return false
}

func (c Claims) HasScope(s string) bool {
	for _, x := range c.Scope {
		if x == s {
			return true
		}
	}
	return false
}

func (c Claims) Attr(key string) (string, bool) {
	if c.Attrs == nil {
		return "", false
	}
	v, ok := c.Attrs[key]
	return v, ok
}
