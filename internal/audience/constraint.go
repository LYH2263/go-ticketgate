package audience

type Constraint struct {
	Issuer   string
	Audience []string
	Subject  string
	UseGlob  bool
}

func (c Constraint) Empty() bool {
	return c.Issuer == "" && len(c.Audience) == 0 && c.Subject == ""
}

func CloneStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
