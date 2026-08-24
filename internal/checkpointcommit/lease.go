package checkpointcommit

type Lease struct {
	ID      string
	closeFn func() error
}

func NewLease(id string, closeFn func() error) *Lease {
	return &Lease{ID: id, closeFn: closeFn}
}

func (l *Lease) Release() error {
	if l == nil || l.closeFn == nil {
		return nil
	}
	_ = l.closeFn()
	return nil
}
