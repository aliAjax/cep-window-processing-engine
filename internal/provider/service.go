package provider

import "fmt"

type Service struct {
	Provider Provider
}

func (s Service) PutVariable(key, value string) error {
	if err := s.Provider.Put(key, value); err != nil {
		return fmt.Errorf("store pattern variable: %w", err)
	}
	return nil
}
