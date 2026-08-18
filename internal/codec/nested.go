package codec

import "fmt"

// EncodeAttrs 将字符串字典编成嵌套块：count u16 + 重复 (kLen u8, k, vLen u16, v)。
func EncodeAttrs(attrs map[string]string) ([]byte, error) {
	if len(attrs) == 0 {
		return nil, nil
	}
	if len(attrs) > MaxAttrs {
		return nil, errTooLarge
	}
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		j := i
		for j > 0 && keys[j] < keys[j-1] {
			keys[j], keys[j-1] = keys[j-1], keys[j]
			j--
		}
	}
	w := NewWriter(MaxPayloadBytes)
	w.U16(uint16(len(keys)))
	for _, k := range keys {
		if len(k) == 0 || len(k) > MaxAttrKey {
			return nil, fmt.Errorf("codec: attr key length")
		}
		v := attrs[k]
		if len(v) > MaxAttrVal {
			return nil, errTooLarge
		}
		if len(k) > 255 {
			return nil, errTooLarge
		}
		w.U8(byte(len(k)))
		w.Raw([]byte(k))
		if OverflowU16(len(v)) {
			return nil, errTooLarge
		}
		w.U16(uint16(len(v)))
		w.Raw([]byte(v))
	}
	if w.Err() != nil {
		return nil, w.Err()
	}
	return w.Bytes(), nil
}

func DecodeAttrs(b []byte) (map[string]string, error) {
	if len(b) == 0 {
		return nil, nil
	}
	r := NewReader(b)
	n, err := r.U16()
	if err != nil {
		return nil, err
	}
	if int(n) > MaxAttrs {
		return nil, errTooLarge
	}
	out := make(map[string]string, n)
	var last string
	for i := 0; i < int(n); i++ {
		kl, err := r.U8()
		if err != nil {
			return nil, err
		}
		kb, err := r.Slice(int(kl))
		if err != nil {
			return nil, err
		}
		vl, err := r.U16()
		if err != nil {
			return nil, err
		}
		vb, err := r.Slice(int(vl))
		if err != nil {
			return nil, err
		}
		k := string(kb)
		if last != "" && k < last {
			return nil, fmt.Errorf("codec: attr key order")
		}
		if _, ok := out[k]; ok {
			return nil, fmt.Errorf("codec: duplicate attr %s", k)
		}
		out[k] = string(vb)
		last = k
	}
	if r.Remaining() != 0 {
		return nil, fmt.Errorf("codec: attr trailing bytes")
	}
	return out, nil
}
