package bind

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(parts ...[]byte) []byte {
	h := sha256.New()
	for _, p := range parts {
		var ln [2]byte
		n := len(p)
		ln[0] = byte(n >> 8)
		ln[1] = byte(n)
		_, _ = h.Write(ln[:])
		_, _ = h.Write(p)
	}
	sum := h.Sum(nil)
	out := make([]byte, 16)
	copy(out, sum[:16])
	return out
}

func HashString(s string) []byte {
	return Hash([]byte(s))
}

func Equal(a, b []byte) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return hmac.Equal(a, b)
}

func Hex(b []byte) string { return hex.EncodeToString(b) }
