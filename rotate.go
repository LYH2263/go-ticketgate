package ticketgate

import "time"

// RotateTicket 验证旧票、签发新票并吊销旧 jti（票据轮换，非密钥轮换）。
func (g *Gateway) RotateTicket(old []byte, extra Claims) (Issued, error) {
	res, err := g.VerifyDetailed(old)
	if err != nil {
		return Issued{}, err
	}
	base := res.Claims.Clone()
	if extra.Subject != "" {
		base.Subject = extra.Subject
	}
	if extra.Issuer != "" {
		base.Issuer = extra.Issuer
	}
	if len(extra.Audience) > 0 {
		base.Audience = append([]string(nil), extra.Audience...)
	}
	if len(extra.Scope) > 0 {
		base.Scope = append([]string(nil), extra.Scope...)
	}
	if extra.Attrs != nil {
		if base.Attrs == nil {
			base.Attrs = map[string]string{}
		}
		for k, v := range extra.Attrs {
			base.Attrs[k] = v
		}
	}
	if len(extra.Bind) > 0 {
		base.Bind = append([]byte(nil), extra.Bind...)
	}
	base.ID = extra.ID
	base.Nonce = extra.Nonce
	base.IssuedAt = timeZero()
	base.NotBefore = extra.NotBefore
	base.ExpiresAt = extra.ExpiresAt
	issued, err := g.IssueDetailed(base)
	if err != nil {
		return Issued{}, err
	}
	if err := g.Revoke(res.Claims.ID); err != nil {
		return Issued{}, err
	}
	return issued, nil
}

func timeZero() time.Time { return time.Time{} }
