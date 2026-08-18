package bind

func Concat(parts ...string) []byte {
	var chunks [][]byte
	for _, p := range parts {
		chunks = append(chunks, []byte(p))
	}
	return Hash(chunks...)
}

func MatchOptional(expect, got []byte) bool {
	if len(expect) == 0 {
		return true
	}
	return Equal(expect, got)
}
