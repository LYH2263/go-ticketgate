package audience

func MatchGlob(pattern, s string) bool {
	return glob(pattern, s)
}

func glob(p, s string) bool {
	for {
		if p == "" {
			return s == ""
		}
		switch p[0] {
		case '*':
			for len(p) > 0 && p[0] == '*' {
				p = p[1:]
			}
			if p == "" {
				return true
			}
			for i := 0; i <= len(s); i++ {
				if glob(p, s[i:]) {
					return true
				}
			}
			return false
		case '?':
			if s == "" {
				return false
			}
			p = p[1:]
			s = s[1:]
		default:
			if s == "" || p[0] != s[0] {
				return false
			}
			p = p[1:]
			s = s[1:]
		}
	}
}
