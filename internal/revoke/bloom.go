package revoke

import "sync"

type Bloom struct {
	mu   sync.Mutex
	m    uint64
	k    int
	bits []uint64
}

func NewBloom(bits uint64, hashes int) *Bloom {
	if bits < 64 {
		bits = 1024
	}
	if hashes < 1 {
		hashes = 4
	}
	nwords := (bits + 63) / 64
	return &Bloom{m: bits, k: hashes, bits: make([]uint64, nwords)}
}

func (b *Bloom) Add(key []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	h1 := FNV64a(key)
	h2 := Mix(h1, 0x9e3779b97f4a7c15)
	for i := 0; i < b.k; i++ {
		idx := (h1 + uint64(i)*h2) % b.m
		b.set(idx)
	}
}

func (b *Bloom) MightContain(key []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	h1 := FNV64a(key)
	h2 := Mix(h1, 0x9e3779b97f4a7c15)
	for i := 0; i < b.k; i++ {
		idx := (h1 + uint64(i)*h2) % b.m
		if !b.get(idx) {
			return false
		}
	}
	return true
}

func (b *Bloom) set(i uint64) {
	w := i / 64
	bit := i % 64
	b.bits[w] |= 1 << bit
}

func (b *Bloom) get(i uint64) bool {
	w := i / 64
	bit := i % 64
	return b.bits[w]&(1<<bit) != 0
}

func (b *Bloom) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.bits {
		b.bits[i] = 0
	}
}
