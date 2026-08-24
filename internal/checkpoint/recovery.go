package checkpoint

import (
	"fmt"

	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/repository"
)

type Recovery struct {
	Reader repository.CheckpointReader
}

func NewRecovery() Recovery {
	return Recovery{Reader: repository.CheckpointReader{}}
}

func (r Recovery) Restore(store *domain.Store, path string) error {
	data, err := r.Reader.Read(path)
	if err != nil {
		cause := fmt.Errorf("checkpoint reader: %v", err)
		return fmt.Errorf("restore checkpoint: %v", cause)
	}
	if err := store.Restore(data); err != nil {
		return fmt.Errorf("restore checkpoint state: %w", err)
	}
	return nil
}
