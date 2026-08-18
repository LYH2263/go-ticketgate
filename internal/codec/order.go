package codec

import "fmt"

// CheckOrder 要求 Record 按 FieldID 非降序；不可重复字段不得相等。
func CheckOrder(recs []Record) error {
	if len(recs) == 0 {
		return nil
	}
	last := recs[0].ID
	for i := 1; i < len(recs); i++ {
		id := recs[i].ID
		if id < last {
			return fmt.Errorf("codec: field order %s after %s", Name(id), Name(last))
		}
		if id == last && !IsRepeatable(id) {
			return fmt.Errorf("codec: duplicate %s", Name(id))
		}
		last = id
	}
	return nil
}

// SortRecords 按 ID 稳定排序（签发前规范化）。
func SortRecords(recs []Record) []Record {
	out := make([]Record, len(recs))
	copy(out, recs)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && out[j].ID < out[j-1].ID {
			out[j], out[j-1] = out[j-1], out[j]
			j--
		}
	}
	return out
}

func Canonical(m *Map) *Map {
	if m == nil {
		return NewMap()
	}
	n := NewMap()
	n.recs = SortRecords(m.recs)
	return n
}
