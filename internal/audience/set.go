package audience

func Normalize(list []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(list))
	for _, s := range list {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func Overlap(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(a))
	for _, x := range a {
		set[x] = struct{}{}
	}
	var out []string
	for _, x := range b {
		if _, ok := set[x]; ok {
			out = append(out, x)
		}
	}
	return out
}
