package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

func newMAC(alg Alg, key []byte) hash.Hash {
	switch alg {
	case AlgHS256:
		return hmac.New(sha256.New, key)
	case AlgHS384:
		return hmac.New(sha512.New384, key)
	case AlgHS512:
		return hmac.New(sha512.New, key)
	default:
		return nil
	}
}

func ComputeMAC(alg Alg, key, message []byte) ([]byte, error) {
	h := newMAC(alg, key)
	if h == nil {
		return nil, errAlg
	}
	if _, err := h.Write(message); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func VerifyMAC(alg Alg, key, message, got []byte) bool {
	want, err := ComputeMAC(alg, key, message)
	if err != nil {
		return false
	}
	if len(want) == 0 || len(got) != len(want) {
		return false
	}
	return hmac.Equal(want, got)
}

var errAlg = errString("token: algorithm")

type errString string

func (e errString) Error() string { return string(e) }
