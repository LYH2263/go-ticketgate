package scope

func Missing(have, need []string) []string {
	set := make(map[string]struct{}, len(have))
	for _, h := range have {
		set[h] = struct{}{}
	}
	var miss []string
	for _, n := range need {
		if _, ok := set[n]; !ok {
			miss = append(miss, n)
		}
	}
	return miss
}

func Has(have []string, s string) bool {
	for _, h := range have {
		if h == s {
			return true
		}
	}
	return false
}
