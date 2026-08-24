package schemaactivation

import (
	"errors"
	"fmt"
)

var ErrSchemaRejected = errors.New("schema rejected")

type RepositoryError struct {
	SchemaID string
	Err      error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("resolve schema %q: %v", e.SchemaID, e.Err)
}
