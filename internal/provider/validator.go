package provider

import (
	"fmt"
)

type Validator struct{}

func (Validator) Validate(value Provider) error {
	if value == nil {
		return fmt.Errorf("provider is required")
	}
	return nil
}
