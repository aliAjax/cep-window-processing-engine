package admission

import "fmt"

type PermanentError struct {
	Err error
}

func (e PermanentError) Error() string { return e.Err.Error() }
func (e PermanentError) Unwrap() error { return e.Err }

type Service struct {
	Repository *Repository
}

func (s Service) Admit(schemaID string) error {
	if s.Repository == nil {
		return fmt.Errorf("admission repository is unavailable")
	}
	_, err := s.Repository.ResolveSchema(schemaID)
	if err != nil {
		return fmt.Errorf("admit pattern: %w", err)
	}
	return nil
}
