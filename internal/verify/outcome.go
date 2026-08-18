package verify

import (
	"time"

	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Stage string

const (
	StageParse   Stage = "parse"
	StageKID     Stage = "kid"
	StageMAC     Stage = "mac"
	StageContext Stage = "context"
	StageTime    Stage = "time"
	StageAud     Stage = "audience"
	StageRevoke  Stage = "revoke"
	StageNonce   Stage = "nonce"
	StageScope   Stage = "scope"
	StageBind    Stage = "bind"
	StageOK      Stage = "ok"
)

type Outcome struct {
	Stage   Stage
	Header  token.Header
	Payload token.Payload
	Now     time.Time
	Err     error
}

func fail(stage Stage, err error) Outcome {
	return Outcome{Stage: stage, Err: err}
}
