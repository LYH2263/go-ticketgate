package revoke

func (r *Revoker) ExactLen() int {
	if r.exact == nil {
		return 0
	}
	return r.exact.Len()
}

func (r *Revoker) Forget(jti string) {
	if r.exact != nil {
		r.exact.Forget(jti)
	}
}

func EstimateBloomFP(n, m uint64, k int) float64 {
	if m == 0 || k == 0 {
		return 1
	}
	// (1 - e^{-kn/m})^k 的粗略展开，避免 math 也可，这里用迭代近似 e^{-x}。
	x := float64(uint64(k)*n) / float64(m)
	ex := approxExpNeg(x)
	p := 1 - ex
	out := 1.0
	for i := 0; i < k; i++ {
		out *= p
	}
	return out
}

func approxExpNeg(x float64) float64 {
	if x < 0 {
		x = 0
	}
	term := 1.0
	sum := 1.0
	for i := 1; i <= 12; i++ {
		term *= -x / float64(i)
		sum += term
	}
	if sum < 0 {
		return 0
	}
	return sum
}
