package token

import (
	"errors"
	"fmt"

	"github.com/LYH2263/go-ticketgate/internal/codec"
)

var ErrCorrupt = errors.New("token: corrupt")

type Envelope struct {
	Version byte
	Header  []byte
	Payload []byte
	MAC     []byte
	Alg     Alg
}

func Encode(h Header, p Payload, secret []byte) ([]byte, error) {
	hb, err := EncodeHeader(h)
	if err != nil {
		return nil, err
	}
	pb, err := EncodePayload(p)
	if err != nil {
		return nil, err
	}
	if len(hb) > codec.MaxHeaderBytes || len(pb) > codec.MaxPayloadBytes {
		return nil, codec.ErrTooLarge()
	}
	macLen := h.Alg.MACSize()
	total := codec.EnvelopeOverhead(macLen) + len(hb) + len(pb)
	if total > codec.MaxTokenBytes {
		return nil, codec.ErrTooLarge()
	}
	w := codec.NewWriter(codec.MaxTokenBytes)
	w.Raw(codec.Magic[:])
	w.U8(codec.VersionV1)
	w.U16(uint16(len(hb)))
	w.U16(uint16(len(pb)))
	w.Raw(hb)
	w.Raw(pb)
	if w.Err() != nil {
		return nil, w.Err()
	}
	prefix := w.Bytes()
	mac, err := ComputeMAC(h.Alg, secret, prefix)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(prefix)+len(mac))
	out = append(out, prefix...)
	out = append(out, mac...)
	return out, nil
}

func Split(raw []byte) (Envelope, error) {
	var env Envelope
	if len(raw) < codec.MinEnvelopeSize(32) {
		return env, fmt.Errorf("%w: truncated", ErrCorrupt)
	}
	if !codec.MagicOK(raw) {
		return env, fmt.Errorf("%w: bad magic", ErrCorrupt)
	}
	r := codec.NewReader(raw)
	if _, err := r.Slice(4); err != nil {
		return env, err
	}
	ver, err := r.U8()
	if err != nil {
		return env, err
	}
	if ver != codec.VersionV1 {
		return env, fmt.Errorf("token: unsupported version %d", ver)
	}
	hl, err := r.U16()
	if err != nil {
		return env, err
	}
	pl, err := r.U16()
	if err != nil {
		return env, err
	}
	if int(hl) > codec.MaxHeaderBytes || int(pl) > codec.MaxPayloadBytes {
		return env, codec.ErrTooLarge()
	}
	hb, err := r.Slice(int(hl))
	if err != nil {
		return env, err
	}
	pb, err := r.Slice(int(pl))
	if err != nil {
		return env, err
	}
	mac := r.Rest()
	if len(mac) != 32 && len(mac) != 48 && len(mac) != 64 {
		return env, fmt.Errorf("token: mac length")
	}
	env.Version = ver
	env.Header = hb
	env.Payload = pb
	env.MAC = mac
	return env, nil
}

func MACMessage(raw []byte, macLen int) ([]byte, error) {
	if len(raw) < codec.EnvelopeOverhead(macLen) {
		return nil, fmt.Errorf("token: truncated")
	}
	return raw[:len(raw)-macLen], nil
}
