package codec

func HeaderPrefixLen() int {
	return 4 + 1 + 2 + 2
}

func SplitSizes(raw []byte) (hdr, pay, mac int, err error) {
	if len(raw) < HeaderPrefixLen()+32 {
		return 0, 0, 0, errTruncated
	}
	if !MagicOK(raw) {
		return 0, 0, 0, errEmpty
	}
	r := NewReader(raw)
	if _, err := r.Slice(4); err != nil {
		return 0, 0, 0, err
	}
	if _, err := r.U8(); err != nil {
		return 0, 0, 0, err
	}
	hl, err := r.U16()
	if err != nil {
		return 0, 0, 0, err
	}
	pl, err := r.U16()
	if err != nil {
		return 0, 0, 0, err
	}
	rest := r.Remaining()
	mac = rest - int(hl) - int(pl)
	if mac < 0 {
		return 0, 0, 0, errTruncated
	}
	return int(hl), int(pl), mac, nil
}

func CopyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
