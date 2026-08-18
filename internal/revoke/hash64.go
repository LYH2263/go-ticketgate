package revoke

const (
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
)

func FNV64a(b []byte) uint64 {
	h := uint64(fnvOffset64)
	for _, c := range b {
		h ^= uint64(c)
		h *= fnvPrime64
	}
	return h
}

func FNV64aString(s string) uint64 {
	return FNV64a([]byte(s))
}

func Mix(h uint64, seed uint64) uint64 {
	h ^= seed
	h *= fnvPrime64
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	return h
}
