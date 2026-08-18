package skew

import "time"

// Window 时钟偏移窗：nbf 允许提前 skew，exp 允许滞后 skew。
// 有效当 (nbf - skew) <= now <= (exp + skew)。边界 now == exp+skew 仍有效（用 After）。
type Window struct {
	Skew time.Duration
}

func New(d time.Duration) Window {
	if d < 0 {
		d = 0
	}
	return Window{Skew: d}
}

func (w Window) Duration() time.Duration { return w.Skew }
