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
	defer func() {
		if releaseErr := guard.Release(); releaseErr != nil && err == nil {
			err = fmt.Errorf("release checkpoint guard: %w", releaseErr)
		}
	}()
	return s.Save(checkpointID)
}
