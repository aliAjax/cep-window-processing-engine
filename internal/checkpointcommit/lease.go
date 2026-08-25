package checkpointcommit

import "fmt"

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
	if err := l.closeFn(); err != nil {
		return fmt.Errorf("release checkpoint lease %q: %w", l.ID, err)
	}
	return nil
}
