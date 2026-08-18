package skew

import "time"

func (w Window) Expand(nbf, exp time.Time) (open, close time.Time) {
	open = nbf
	close = exp
	if !nbf.IsZero() {
		open = nbf.Add(-w.Skew)
	}
	if !exp.IsZero() {
		close = exp.Add(w.Skew)
	}
	return open, close
}

func (w Window) Contains(now, nbf, exp time.Time) bool {
	expired, early := w.Check(now, nbf, exp)
	return !expired && !early
}

func Seconds(d time.Duration) int64 {
	return int64(d / time.Second)
}
