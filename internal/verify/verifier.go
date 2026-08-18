package verify

type Verifier struct {
	deps Deps
}

func New(d Deps) *Verifier {
	return &Verifier{deps: d}
}

func (v *Verifier) Verify(token []byte) (Outcome, error) {
	o := Run(token, v.deps)
	if o.Err != nil {
		return o, o.Err
	}
	return o, nil
}

func (v *Verifier) Deps() Deps { return v.deps }
