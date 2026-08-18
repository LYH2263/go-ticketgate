package keyring

import (
	"fmt"
	"sync"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/clock"
	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Key struct {
	Material Material
	Status   Status
	Created  time.Time
	RetireAt time.Time
}

type Ring struct {
	mu      sync.Mutex
	clk     clock.Clock
	grace   time.Duration
	current string
	keys    map[string]*Key
	order   []string
	closed  bool
}

func New(clk clock.Clock, grace time.Duration, alg token.Alg) (*Ring, error) {
	if clk == nil {
		clk = clock.Real{}
	}
	if grace < 0 {
		grace = 0
	}
	if !alg.Valid() {
		alg = token.AlgHS256
	}
	r := &Ring{
		clk:   clk,
		grace: grace,
		keys:  make(map[string]*Key),
	}
	if _, err := r.rotateLocked(alg); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Ring) GracePeriod() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.grace
}

func (r *Ring) SetGrace(d time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d < 0 {
		d = 0
	}
	r.grace = d
}

func (r *Ring) Current() (Key, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	k, ok := r.keys[r.current]
	if !ok || k.Status != StatusActive {
		return Key{}, fmt.Errorf("keyring: no current key")
	}
	return cloneKey(k), nil
}

func (r *Ring) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
}

func (r *Ring) Get(kid string) (Key, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return Key{}, fmt.Errorf("keyring: closed")
	}
	r.sweepLocked()
	k, ok := r.keys[kid]
	if !ok {
		return Key{}, fmt.Errorf("keyring: unknown kid")
	}
	return cloneKey(k), nil
}

func (r *Ring) LookupVerify(kid string) (Key, error) {
	k, err := r.Get(kid)
	if err != nil {
		return Key{}, err
	}
	if !k.Status.CanVerify() {
		return Key{}, fmt.Errorf("keyring: key retired")
	}
	return k, nil
}

func (r *Ring) List() []Key {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sweepLocked()
	out := make([]Key, 0, len(r.order))
	for _, id := range r.order {
		if k, ok := r.keys[id]; ok {
			out = append(out, cloneKey(k))
		}
	}
	return out
}

func cloneKey(k *Key) Key {
	return Key{
		Material: Material{
			KID:    k.Material.KID,
			Alg:    k.Material.Alg,
			Secret: CloneSecret(k.Material.Secret),
		},
		Status:   k.Status,
		Created:  k.Created,
		RetireAt: k.RetireAt,
	}
}
