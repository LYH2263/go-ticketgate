package scope

func ContainsAll(have, need []string) bool {
	if len(need) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(have))
	for _, h := range have {
		set[h] = struct{}{}
	}
	for _, n := range need {
		if _, ok := set[n]; !ok {
			return false
		}
	}
	return true
}

func Intersect(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(a))
	for _, x := range a {
		set[x] = struct{}{}
	}
	var out []string
	seen := map[string]struct{}{}
	for _, x := range b {
		if _, ok := set[x]; ok {
			if _, d := seen[x]; d {
				continue
			}
			seen[x] = struct{}{}
			out = append(out, x)
		}
	}
	return out
}

func Unique(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
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
