package codec

import (
	"errors"
	"fmt"
)

var ErrCorrupt = errors.New("codec: corrupt")

// Record 单条 TLV。
type Record struct {
	ID    FieldID
	Value []byte
}

type Map struct {
	recs []Record
}

func NewMap() *Map { return &Map{} }

func (m *Map) Add(id FieldID, v []byte) {
	cp := make([]byte, len(v))
	copy(cp, v)
	m.recs = append(m.recs, Record{ID: id, Value: cp})
}

func (m *Map) AddString(id FieldID, s string) { m.Add(id, []byte(s)) }

func (m *Map) Records() []Record { return m.recs }

func (m *Map) First(id FieldID) ([]byte, bool) {
	for _, r := range m.recs {
		if r.ID == id {
			return r.Value, true
		}
	}
	return nil, false
}

func (m *Map) All(id FieldID) [][]byte {
	var out [][]byte
	for _, r := range m.recs {
		if r.ID == id {
			out = append(out, r.Value)
		}
	}
	return out
}

func (m *Map) FirstString(id FieldID) string {
	v, ok := m.First(id)
	if !ok {
		return ""
	}
	return string(v)
}

func (m *Map) AllString(id FieldID) []string {
	raw := m.All(id)
	out := make([]string, 0, len(raw))
	for _, b := range raw {
		out = append(out, string(b))
	}
	return out
}

func EncodeMap(m *Map, max int) ([]byte, error) {
	if m == nil {
		return []byte{}, nil
	}
	if err := CheckOrder(m.recs); err != nil {
		return nil, err
	}
	w := NewWriter(max)
	for _, r := range m.recs {
		spec, ok := Spec(r.ID)
		if !ok {
			return nil, fmt.Errorf("codec: unknown field %s: %w", Name(r.ID), errUnknownField())
		}
		_ = spec
		if OverflowU16(len(r.Value)) {
			return nil, errTooLarge
		}
		w.U8(byte(r.ID))
		w.U16(uint16(len(r.Value)))
		w.Raw(r.Value)
	}
	if w.Err() != nil {
		return nil, w.Err()
	}
	return w.Bytes(), nil
}

func DecodeMap(b []byte, header bool, strictUnknown bool) (*Map, error) {
	r := NewReader(b)
	m := NewMap()
	var last FieldID
	first := true
	seen := make(map[FieldID]int)
	for r.Remaining() > 0 {
		id, val, err := r.TLV()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
		}
		if !first && id < last {
			return nil, errFieldOrder()
		}
		if !first && id == last && !IsRepeatable(id) {
			return nil, errDuplicate()
		}
		spec, ok := Spec(id)
		if !ok {
			if strictUnknown {
				return nil, errUnknownField()
			}
			last = id
			first = false
			continue
		}
		if spec.Header != header {
			return nil, errUnknownField()
		}
		seen[id]++
		if seen[id] > 1 && !spec.Repeatable {
			return nil, errDuplicate()
		}
		m.Add(id, val)
		last = id
		first = false
	}
	return m, nil
}

func errUnknownField() error { return fmt.Errorf("codec: unknown field") }
func errFieldOrder() error   { return fmt.Errorf("codec: field order") }
func errDuplicate() error    { return fmt.Errorf("codec: duplicate field") }

func WrapUnknown() error { return errUnknownField() }
func WrapOrder() error   { return errFieldOrder() }
func WrapDup() error     { return errDuplicate() }
