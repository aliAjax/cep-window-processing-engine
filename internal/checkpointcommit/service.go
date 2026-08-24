package checkpointcommit

import "fmt"

type Guard interface {
	Release() error
}

type Service struct {
	Acquire func(checkpointID string) (Guard, error)
	Save    func(checkpointID string) error
}

func (s Service) Commit(checkpointID string) (err error) {
	guard, err := s.Acquire(checkpointID)
	if err != nil {
		return err
	}
	if err = s.Save(checkpointID); err != nil {
		return err
	}
	if err = guard.Release(); err != nil {
		return fmt.Errorf("release checkpoint guard: %w", err)
	}
	return nil
}
