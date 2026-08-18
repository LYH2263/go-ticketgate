package verify

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/skew"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

func CheckTime(p token.Payload, now time.Time, w skew.Window) error {
	expired, early := w.Check(now, p.NotBefore, p.ExpiresAt)
	if early {
		return fmt.Errorf("verify: not yet valid")
	}
	if expired {
		return fmt.Errorf("verify: expired")
	}
	return nil
}

func Deadline(p token.Payload, w skew.Window) time.Time {
	if p.ExpiresAt.IsZero() {
		return time.Time{}
	}
	return p.ExpiresAt.Add(w.Skew)
}

func OpensAt(p token.Payload, w skew.Window) time.Time {
	if p.NotBefore.IsZero() {
		return time.Time{}
	}
	return p.NotBefore.Add(-w.Skew)
}
