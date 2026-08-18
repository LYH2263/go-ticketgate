package revoke

type Persist interface {
	Save(jtis []string) error
	Load() ([]string, error)
}

func (r *Revoker) SetPersist(p Persist) { r.persist = p }

func (r *Revoker) listJTIs() []string {
	if r.exact == nil {
		return nil
	}
	return r.exact.List()
}

func (r *Revoker) Reload() error {
	if r.persist == nil {
		return nil
	}
	jtis, err := r.persist.Load()
	if err != nil {
		return err
	}
	if r.exact != nil {
		r.exact.Replace(jtis)
	}
	return nil
}
