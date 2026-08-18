package clock

import "time"

func UnixUTC(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

func MaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func MinTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() || a.Before(b) {
		return a
	}
	return b
}
