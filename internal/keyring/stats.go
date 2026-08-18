package keyring

func (r *Ring) CurrentKID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.current
}

func (r *Ring) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.keys)
}

func (r *Ring) ActiveCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	n := 0
	for _, k := range r.keys {
		if k.Status == StatusActive {
			n++
		}
	}
	return n
}

func (r *Ring) RetiringCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	n := 0
	for _, k := range r.keys {
		if k.Status == StatusRetiring {
			n++
		}
	}
	return n
}
