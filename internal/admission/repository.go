package admission

import (
	"errors"
	"fmt"
	"sync"
)

var ErrSchemaNotFound = errors.New("schema not found")

type Schema struct {
	ID      string
	Version string
}

type Repository struct {
	mu      sync.RWMutex
	schemas map[string]Schema
}

func NewRepository() *Repository {
	return &Repository{schemas: make(map[string]Schema)}
}

func (r *Repository) PutSchema(schema Schema) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.schemas[schema.ID] = schema
}

func (r *Repository) ResolveSchema(id string) (Schema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	schema, ok := r.schemas[id]
	if !ok {
		return Schema{}, fmt.Errorf("resolve schema %q: %w", id, ErrSchemaNotFound)
	}
	return schema, nil
}
