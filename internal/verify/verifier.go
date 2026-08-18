package verify

import "context"

type Verifier struct {
	deps Deps
}

func New(d Deps) *Verifier {
	return &Verifier{deps: d}
}

func (v *Verifier) Verify(token []byte) (Outcome, error) {
	return v.VerifyContext(context.Background(), token)
}

func (v *Verifier) VerifyContext(ctx context.Context, token []byte) (Outcome, error) {
	o := RunContext(ctx, token, v.deps)
	if o.Err != nil {
		return o, o.Err
	}
	return o, nil
}

func (v *Verifier) Deps() Deps { return v.deps }
