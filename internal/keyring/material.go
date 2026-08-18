package keyring

import (
	"crypto/rand"
	"fmt"

	"github.com/LYH2263/go-ticketgate/internal/token"
)

type Material struct {
	KID    string
	Alg    token.Alg
	Secret []byte
}

func Generate(alg token.Alg) (Material, error) {
	var m Material
	if !alg.Valid() {
		return m, fmt.Errorf("keyring: alg")
	}
	sec := make([]byte, alg.KeySize())
	if _, err := rand.Read(sec); err != nil {
		return m, err
	}
	kid, err := newKID()
	if err != nil {
		return m, err
	}
	m.KID = kid
	m.Alg = alg
	m.Secret = sec
	return m, nil
}

func newKID() (string, error) {
	var b [6]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("k-%x", b[:]), nil
}

func CloneSecret(s []byte) []byte {
	out := make([]byte, len(s))
	copy(out, s)
	return out
}

func ValidMaterial(m Material) error {
	if m.KID == "" {
		return fmt.Errorf("keyring: empty kid")
	}
	if !m.Alg.Valid() {
		return fmt.Errorf("keyring: alg")
	}
	if len(m.Secret) < m.Alg.KeySize() {
		return fmt.Errorf("keyring: short secret")
	}
	return nil
}
