package keyring

func (r *Ring) Sweep() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
}

func (r *Ring) sweepLocked() {
	now := r.clk.Now()
	for _, k := range r.keys {
		if k.Status == StatusRetiring && !k.RetireAt.IsZero() && !now.Before(k.RetireAt) {
			k.Status = StatusRetired
		}
	}
}

// InGrace 报告 kid 是否仍处于宽限期（可验证、不可签发）。
func (r *Ring) InGrace(kid string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	k, ok := r.keys[kid]
	if !ok {
		return false
	}
	return k.Status == StatusRetiring
}

func (r *Ring) Retired(kid string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	k, ok := r.keys[kid]
	if !ok {
		return true
	}
	return k.Status == StatusRetired
}
