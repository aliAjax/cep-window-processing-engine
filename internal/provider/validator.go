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
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return fmt.Errorf("provider is nil")
	}
	return nil
}
