package token

import "fmt"

type Alg byte

const (
	AlgNone  Alg = 0
	AlgHS256 Alg = 1
	AlgHS384 Alg = 2
	AlgHS512 Alg = 3
)

func (a Alg) String() string {
	switch a {
	case AlgHS256:
		return "HS256"
	case AlgHS384:
		return "HS384"
	case AlgHS512:
		return "HS512"
	default:
		return fmt.Sprintf("alg-%d", int(a))
	}
}

func (a Alg) MACSize() int {
	switch a {
	case AlgHS256:
		return 32
	case AlgHS384:
		return 48
	case AlgHS512:
		return 64
	default:
		return 0
	}
}

func (a Alg) KeySize() int {
	switch a {
	case AlgHS256:
		return 32
	case AlgHS384:
		return 48
	case AlgHS512:
		return 64
	default:
		return 0
	}
}

func (a Alg) Valid() bool {
	return a == AlgHS256 || a == AlgHS384 || a == AlgHS512
}

func ParseAlg(b byte) (Alg, bool) {
	a := Alg(b)
	return a, a.Valid()
}
