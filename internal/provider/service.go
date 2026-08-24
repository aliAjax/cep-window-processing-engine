package provider

import "fmt"

type Service struct {
	Provider Provider
}

func (s Service) PutVariable(key, value string) error {
	if s.Provider == nil {
		return fmt.Errorf("provider is required")
	}
	if err := s.Provider.Put(key, value); err != nil {
		return fmt.Errorf("store pattern variable: %w", err)
	}
	return nil
}
