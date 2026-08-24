package schemaactivation

import "fmt"

type PolicyError struct {
	Decision string
	Err      error
}

func (e *PolicyError) Error() string {
	return fmt.Sprintf("schema activation decision %q: %v", e.Decision, e.Err)
}
