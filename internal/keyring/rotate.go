package keyring

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-ticketgate/internal/token"
)

func (r *Ring) Rotate() (string, error) {
	return r.RotateAlg(token.AlgHS256)
}

func (r *Ring) RotateAlg(alg token.Alg) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rotateLocked(alg)
}

func (r *Ring) rotateLocked(alg token.Alg) (string, error) {
	if !alg.Valid() {
		return "", fmt.Errorf("keyring: alg")
	}
	now := r.clk.Now()
	if r.current != "" {
		if old, ok := r.keys[r.current]; ok && old.Status == StatusActive {
			old.Status = StatusRetiring
			old.RetireAt = now.Add(r.grace)
			if r.grace == 0 {
				old.Status = StatusRetired
			}
		}
	}
	mat, err := Generate(alg)
	if err != nil {
		return "", err
	}
	k := &Key{
		Material: mat,
		Status:   StatusActive,
		Created:  now,
	}
	r.keys[mat.KID] = k
	r.order = append(r.order, mat.KID)
	r.current = mat.KID
	r.sweepLocked()
	return mat.KID, nil
}

func (r *Ring) Import(m Material, active bool) error {
	if err := ValidMaterial(m); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.keys[m.KID]; ok {
		return fmt.Errorf("keyring: duplicate kid")
	}
	now := r.clk.Now()
	st := StatusRetiring
	var retireAt = now.Add(r.grace)
	if active {
		if r.current != "" {
			if old, ok := r.keys[r.current]; ok && old.Status == StatusActive {
				old.Status = StatusRetiring
				old.RetireAt = now.Add(r.grace)
			}
		}
		st = StatusActive
		retireAt = time.Time{}
		r.current = m.KID
	}
	r.keys[m.KID] = &Key{
		Material: Material{KID: m.KID, Alg: m.Alg, Secret: CloneSecret(m.Secret)},
		Status:   st,
		Created:  now,
		RetireAt: retireAt,
	}
	r.order = append(r.order, m.KID)
	return nil
}

func (r *Ring) Has(kid string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.keys[kid]
	return ok
}
