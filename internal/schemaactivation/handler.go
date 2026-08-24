package schemaactivation

import "fmt"

type HandlerError struct {
	Status int
	Err    error
}

func (e *HandlerError) Error() string {
	return fmt.Sprintf("schema activation status %d: %v", e.Status, e.Err)
}
