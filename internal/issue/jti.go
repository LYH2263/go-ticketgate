package issue

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func RandomHex(n int) (string, error) {
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewJTI() (string, error) {
	s, err := RandomHex(12)
	if err != nil {
		return "", err
	}
	return "j-" + s, nil
}

func NewNonce() (string, error) {
	s, err := RandomHex(12)
	if err != nil {
		return "", err
	}
	return "n-" + s, nil
}

func MustID(existing string, gen func() (string, error), required bool, what string) (string, error) {
	if existing != "" {
		return existing, nil
	}
	if !required && gen == nil {
		return "", nil
	}
	if gen == nil {
		return "", fmt.Errorf("issue: missing %s", what)
	}
	s, err := gen()
	if err != nil {
		return "", err
	}
	if s == "" && required {
		return "", fmt.Errorf("issue: empty %s", what)
	}
	return s, nil
}
