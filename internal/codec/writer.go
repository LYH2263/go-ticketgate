package codec

import "encoding/binary"

type Writer struct {
	buf []byte
	cap int
	err error
}

func NewWriter(max int) *Writer {
	if max <= 0 {
		max = MaxTokenBytes
	}
	return &Writer{buf: make([]byte, 0, 256), cap: max}
}

func (w *Writer) Err() error { return w.err }

func (w *Writer) Bytes() []byte {
	if w.err != nil {
		return nil
	}
	out := make([]byte, len(w.buf))
	copy(out, w.buf)
	return out
}

func (w *Writer) Len() int { return len(w.buf) }

func (w *Writer) grow(n int) bool {
	if w.err != nil {
		return false
	}
	if len(w.buf)+n > w.cap {
		w.err = errTooLarge
		return false
	}
	return true
}

func (w *Writer) U8(v byte) {
	if !w.grow(1) {
		return
	}
	w.buf = append(w.buf, v)
}

func (w *Writer) U16(v uint16) {
	if !w.grow(2) {
		return
	}
	var tmp [2]byte
	binary.BigEndian.PutUint16(tmp[:], v)
	w.buf = append(w.buf, tmp[:]...)
}

func (w *Writer) U32(v uint32) {
	if !w.grow(4) {
		return
	}
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], v)
	w.buf = append(w.buf, tmp[:]...)
}

func (w *Writer) U64(v uint64) {
	if !w.grow(8) {
		return
	}
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], v)
	w.buf = append(w.buf, tmp[:]...)
}

func (w *Writer) Raw(b []byte) {
	if !w.grow(len(b)) {
		return
	}
	w.buf = append(w.buf, b...)
}

func (w *Writer) TLVBytes(id FieldID, b []byte) {
	if OverflowU16(len(b)) {
		w.err = errTooLarge
		return
	}
	w.U8(byte(id))
	w.U16(uint16(len(b)))
	w.Raw(b)
}

func (w *Writer) TLVString(id FieldID, s string) {
	w.TLVBytes(id, []byte(s))
}

func (w *Writer) TLVU8(id FieldID, v byte) {
	w.U8(byte(id))
	w.U16(1)
	w.U8(v)
}

func (w *Writer) TLVU16(id FieldID, v uint16) {
	w.U8(byte(id))
	w.U16(2)
	w.U16(v)
}

func (w *Writer) TLVU64(id FieldID, v uint64) {
	w.U8(byte(id))
	w.U16(8)
	w.U64(v)
}

func (w *Writer) Reset() {
	w.buf = w.buf[:0]
	w.err = nil
}
