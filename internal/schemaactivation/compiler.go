package schemaactivation

import "fmt"

type CompileError struct {
	Version string
	Err     error
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("compile schema version %q: %v", e.Version, e.Err)
}
