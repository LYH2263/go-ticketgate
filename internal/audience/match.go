package audience

func MatchIssuer(want, got string) bool {
	if want == "" {
		return true
	}
	return want == got
}

func MatchSubject(want, got string) bool {
	if want == "" {
		return true
	}
	return want == got
}

// MatchAudience：期望列表非空时，票据 aud 必须与期望有交集。
// UseGlob 时期望项可作为模式匹配票据 aud。
func MatchAudience(want, got []string, useGlob bool) bool {
	if len(want) == 0 {
		return true
	}
	if len(got) == 0 {
		return false
	}
	for _, w := range want {
		for _, g := range got {
			if useGlob {
				if MatchGlob(w, g) || MatchGlob(g, w) {
					return true
				}
			} else if w == g {
				return true
			}
		}
	}
	return false
}

func (c Constraint) Check(iss, sub string, aud []string) (issOK, subOK, audOK bool) {
	return MatchIssuer(c.Issuer, iss), MatchSubject(c.Subject, sub), MatchAudience(c.Audience, aud, c.UseGlob)
}
