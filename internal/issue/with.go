package issue

import "time"

func (p Policy) WithIssuer(iss string) Policy {
	p.DefaultIssuer = iss
	return p
}

func (p Policy) WithDefaultTTL(d time.Duration) Policy {
	p.DefaultTTL = d
	return p
}

func (p Policy) WithBounds(min, max time.Duration) Policy {
	p.MinTTL = min
	p.MaxTTL = max
	return p
}

func (p Policy) WithDefaultAUD(aud ...string) Policy {
	p.DefaultAUD = append([]string(nil), aud...)
	p.RequireAUD = len(aud) > 0
	return p
}

func (d Draft) WithSubject(sub string) Draft {
	d.Subject = sub
	return d
}

func (d Draft) WithAudience(aud ...string) Draft {
	d.Audience = append([]string(nil), aud...)
	return d
}

func (d Draft) WithTTL(ttl time.Duration) Draft {
	d.TTL = ttl
	return d
}

func (d Draft) WithScope(scope ...string) Draft {
	d.Scope = append([]string(nil), scope...)
	return d
}

func (d Draft) WithAttr(k, v string) Draft {
	if d.Attrs == nil {
		d.Attrs = map[string]string{}
	}
	d.Attrs[k] = v
	return d
}
