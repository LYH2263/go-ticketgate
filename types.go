package ticketgate

import "time"

// Claims 是票据载荷。ID 为 jti；Nonce 用于防重放。
type Claims struct {
	Issuer    string
	Subject   string
	Audience  []string
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	ID        string
	Nonce     string
	Scope     []string
	Attrs     map[string]string
	Bind      []byte
}

// Issued 签发结果，含回填后的声明。
type Issued struct {
	Token  []byte
	KID    string
	Claims Claims
}

// KeyInfo 描述密钥环中的一把密钥。
type KeyInfo struct {
	KID      string
	Alg      byte
	Status   string
	Created  time.Time
	RetireAt time.Time
}

// VerifyResult 验证成功时的附加上下文。
type VerifyResult struct {
	Claims Claims
	KID    string
	Alg    byte
}
