package provider

import (
	"fmt"
	"reflect"
)

type Validator struct{}

func (Validator) Validate(value Provider) error {
	if value == nil {
		return fmt.Errorf("provider is required")
	}
	// Guard against typed-nil providers (a concrete pointer wrapped in a
	// non-nil interface), which would panic on first use.
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("provider is required")
	}
	return nil
}
