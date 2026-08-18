package verify

import "github.com/LYH2263/go-ticketgate/internal/audience"

func (e Expect) WithAudience(aud ...string) Expect {
	e.Audience.Audience = append([]string(nil), aud...)
	return e
}

func (e Expect) WithIssuer(iss string) Expect {
	e.Audience.Issuer = iss
	return e
}

func (e Expect) WithSubject(sub string) Expect {
	e.Audience.Subject = sub
	return e
}

func (e Expect) WithGlob() Expect {
	e.Audience.UseGlob = true
	return e
}

func (e Expect) WithScopes(s ...string) Expect {
	e.Scopes = append([]string(nil), s...)
	return e
}

func (e Expect) Constraint() audience.Constraint {
	return e.Audience
}
