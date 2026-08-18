package token

func BuildFlags(p Payload) uint16 {
	var f uint16
	if len(p.Bind) > 0 {
		f |= FlagHasBind
	}
	if len(p.Scope) > 0 {
		f |= FlagHasScope
	}
	if len(p.Attrs) > 0 {
		f |= FlagHasAttrs
	}
	return f
}

func IssueBytes(kid string, alg Alg, secret []byte, p Payload) ([]byte, error) {
	h := Header{KID: kid, Alg: alg, Flags: BuildFlags(p)}
	return Encode(h, p, secret)
}
