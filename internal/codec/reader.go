package codec

import (
	"encoding/binary"
	"errors"
)

var (
	errTooLarge  = errors.New("codec: too large")
	errTruncated = errors.New("codec: truncated")
	errEmpty     = errors.New("codec: empty")
)

func ErrTooLarge() error  { return errTooLarge }
func ErrTruncated() error { return errTruncated }

type Reader struct {
	b   []byte
	off int
}

func NewReader(b []byte) *Reader {
	return &Reader{b: b}
}

func (r *Reader) Remaining() int {
	if r.off >= len(r.b) {
		return 0
	}
	return len(r.b) - r.off
}

func (r *Reader) Offset() int { return r.off }

func (r *Reader) Rest() []byte {
	if r.off >= len(r.b) {
		return nil
	}
	out := make([]byte, len(r.b)-r.off)
	copy(out, r.b[r.off:])
	return out
}

func (r *Reader) need(n int) error {
	if n < 0 || r.off+n > len(r.b) {
		return errTruncated
	}
	return nil
}

func (r *Reader) U8() (byte, error) {
	if err := r.need(1); err != nil {
		return 0, err
	}
	v := r.b[r.off]
	r.off++
	return v, nil
}

func (r *Reader) U16() (uint16, error) {
	if err := r.need(2); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint16(r.b[r.off:])
	r.off += 2
	return v, nil
}

func (r *Reader) U32() (uint32, error) {
	if err := r.need(4); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(r.b[r.off:])
	r.off += 4
	return v, nil
}

func (r *Reader) U64() (uint64, error) {
	if err := r.need(8); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(r.b[r.off:])
	r.off += 8
	return v, nil
}

func (r *Reader) Slice(n int) ([]byte, error) {
	if err := r.need(n); err != nil {
		return nil, err
	}
	out := make([]byte, n)
	copy(out, r.b[r.off:r.off+n])
	r.off += n
	return out, nil
}

func (r *Reader) Skip(n int) error {
	if err := r.need(n); err != nil {
		return err
	}
	r.off += n
	return nil
}

func (r *Reader) TLV() (FieldID, []byte, error) {
	idb, err := r.U8()
	if err != nil {
		return 0, nil, err
	}
	n, err := r.U16()
	if err != nil {
		return 0, nil, err
	}
	body, err := r.Slice(int(n))
	if err != nil {
		return 0, nil, err
	}
	return FieldID(idb), body, nil
}

func AsU8(b []byte) (byte, error) {
	if len(b) != 1 {
		return 0, errTruncated
	}
	return b[0], nil
}

func AsU16(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, errTruncated
	}
	return binary.BigEndian.Uint16(b), nil
}

func AsU64(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, errTruncated
	}
	return binary.BigEndian.Uint64(b), nil
}

func PutU16(dst []byte, v uint16) {
	binary.BigEndian.PutUint16(dst, v)
}

func PutU64(dst []byte, v uint64) {
	binary.BigEndian.PutUint64(dst, v)
}
