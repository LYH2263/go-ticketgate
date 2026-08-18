package skew

import (
	"context"
	"time"
)

func Expired(now, exp time.Time, skew time.Duration) bool {
	if exp.IsZero() {
		return true
	}
	deadline := exp.Add(skew)
	return now.After(deadline)
}

func TooEarly(now, nbf time.Time, skew time.Duration) bool {
	if nbf.IsZero() {
		return false
	}
	open := nbf.Add(-skew)
	return now.Before(open)
}

func (w Window) Check(now, nbf, exp time.Time) (expired, early bool) {
	return Expired(now, exp, w.Skew), TooEarly(now, nbf, w.Skew)
}

func (w Window) CheckContext(ctx context.Context, now, nbf, exp time.Time) (expired, early bool, err error) {
	_ = ctx
	return Expired(now, exp, w.Skew), TooEarly(now, nbf, w.Skew), nil
}

func InstantUnix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

func EqualSec(a, b time.Time) bool {
	return InstantUnix(a) == InstantUnix(b)
}
