package token

import (
	"fmt"

	"github.com/LYH2263/go-ticketgate/internal/codec"
)

const (
	FlagNone     uint16 = 0
	FlagHasBind  uint16 = 1 << 0
	FlagHasScope uint16 = 1 << 1
	FlagHasAttrs uint16 = 1 << 2
)

type Header struct {
	KID   string
	Alg   Alg
	Flags uint16
}

func EncodeHeader(h Header) ([]byte, error) {
	if h.KID == "" || len(h.KID) > codec.MaxKidBytes {
		return nil, fmt.Errorf("token: kid")
	}
	if !h.Alg.Valid() {
		return nil, errAlg
	}
	m := codec.NewMap()
	m.AddString(codec.FKid, h.KID)
	m.Add(codec.FAlg, []byte{byte(h.Alg)})
	var fl [2]byte
	codec.PutU16(fl[:], h.Flags)
	m.Add(codec.FFlags, fl[:])
	return codec.EncodeMap(m, codec.MaxHeaderBytes)
}

func DecodeHeader(b []byte) (Header, error) {
	var h Header
	m, err := codec.DecodeMap(b, true, true)
	if err != nil {
		return h, err
	}
	kid := m.FirstString(codec.FKid)
	if kid == "" || len(kid) > codec.MaxKidBytes {
		return h, fmt.Errorf("token: kid")
	}
	ab, ok := m.First(codec.FAlg)
	if !ok {
		return h, fmt.Errorf("token: alg missing")
	}
	algB, err := codec.AsU8(ab)
	if err != nil {
		return h, err
	}
	alg, ok := ParseAlg(algB)
	if !ok {
		return h, errAlg
	}
	var flags uint16
	if fb, ok := m.First(codec.FFlags); ok {
		flags, err = codec.AsU16(fb)
		if err != nil {
			return h, err
		}
	}
	h.KID = kid
	h.Alg = alg
	h.Flags = flags
	return h, nil
}
