package verify

import (
	"fmt"

	"github.com/LYH2263/go-ticketgate/internal/audience"
	"github.com/LYH2263/go-ticketgate/internal/bind"
	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/keyring"
	"github.com/LYH2263/go-ticketgate/internal/nonce"
	"github.com/LYH2263/go-ticketgate/internal/revoke"
	"github.com/LYH2263/go-ticketgate/internal/scope"
	"github.com/LYH2263/go-ticketgate/internal/skew"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Expect struct {
	Audience audience.Constraint
	Scopes   []string
	Bind     []byte
}

type Deps struct {
	Ring    *keyring.Ring
	Revoker *revoke.Revoker
	Nonces  *nonce.Cache
	Skew    skew.Window
	Clock   clock.Clock
	Expect  Expect
}

func Run(raw []byte, d Deps) Outcome {
	if d.Clock == nil {
		d.Clock = clock.Real{}
	}
	now := d.Clock.Now()

	env, err := token.Split(raw)
	if err != nil {
		return fail(StageParse, err)
	}
	h, err := token.DecodeHeader(env.Header)
	if err != nil {
		return fail(StageParse, err)
	}

	key, err := d.Ring.LookupVerify(h.KID)
	if err != nil {
		st := StageKID
		if err.Error() == "keyring: key retired" {
			st = StageKID
		}
		return fail(st, err)
	}
	if h.Alg != key.Material.Alg {
		return fail(StageMAC, fmt.Errorf("verify: alg mismatch"))
	}
	if h.Alg.MACSize() != len(env.MAC) {
		return fail(StageMAC, fmt.Errorf("verify: mac length"))
	}
	msg, err := token.MACMessage(raw, len(env.MAC))
	if err != nil {
		return fail(StageMAC, err)
	}
	if !token.VerifyMAC(h.Alg, key.Material.Secret, msg, env.MAC) {
		return fail(StageMAC, fmt.Errorf("verify: bad mac"))
	}

	p, err := token.DecodePayload(env.Payload)
	if err != nil {
		return fail(StageParse, err)
	}

	if err := CheckTime(p, now, d.Skew); err != nil {
		return Outcome{Stage: StageTime, Header: h, Payload: p, Now: now, Err: err}
	}

	issOK, subOK, audOK := d.Expect.Audience.Check(p.Issuer, p.Subject, p.Audience)
	if !issOK {
		return Outcome{Stage: StageAud, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: issuer mismatch")}
	}
	if !subOK {
		return Outcome{Stage: StageAud, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: subject mismatch")}
	}
	if !audOK {
		return Outcome{Stage: StageAud, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: audience mismatch")}
	}

	if d.Revoker != nil && p.JTI != "" && d.Revoker.IsRevoked(p.JTI) {
		return Outcome{Stage: StageRevoke, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: revoked")}
	}

	if d.Nonces != nil && p.Nonce != "" {
		until := p.ExpiresAt.Add(d.Skew.Skew)
		if !d.Nonces.Occupy(p.Nonce, until) {
			return Outcome{Stage: StageNonce, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: nonce replay")}
		}
	}

	if !scope.ContainsAll(p.Scope, d.Expect.Scopes) {
		return Outcome{Stage: StageScope, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: scope")}
	}
	if len(d.Expect.Bind) > 0 && !bind.Equal(d.Expect.Bind, p.Bind) {
		return Outcome{Stage: StageBind, Header: h, Payload: p, Now: now, Err: fmt.Errorf("verify: bind mismatch")}
	}

	return Outcome{Stage: StageOK, Header: h, Payload: p, Now: now}
}
