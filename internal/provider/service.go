package provider

import "fmt"

type Service struct {
	Provider Provider
}

func (s Service) PutVariable(key, value string) error {
	if err := (Validator{}).Validate(s.Provider); err != nil {
		return fmt.Errorf("put pattern variable: %w", err)
	}
	if err := s.Provider.Put(key, value); err != nil {
		return fmt.Errorf("store pattern variable: %w", err)
	}
	return nil
}
