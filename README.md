# go-ticketgate

会话票据网关库：签发、验证、轮换、吊销短期票据。自定义紧凑二进制格式（非 JWT），支持密钥轮换与宽限期、时钟偏移窗、受众/签发者约束、nonce 防重放，以及精确/布隆吊销。无前端。

## 不变量

1. 过期票验证失败；偏移窗内边界可测。
2. 吊销后验证失败。
3. 旧密钥在宽限期内仍可验；宽限外失败。
4. 受众不匹配失败。
5. 同 nonce 重放失败（在缓存窗内）。

## 用法

```go
g := ticketgate.New(
    ticketgate.WithIssuer("my-app"),
    ticketgate.WithAudience("api"),
    ticketgate.WithTTL(15*time.Minute),
    ticketgate.WithGrace(1*time.Hour),
)
tok, kid, err := g.Issue(ticketgate.Claims{Subject: "user-1"})
claims, err := g.Verify(tok)
_ = kid
_ = claims
g.Revoke(claims.ID)
g.RotateKey()
```

## 测试

```bash
go test ./... -count=1
```
