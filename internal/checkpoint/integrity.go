package checkpoint

import (
	"github.com/example/cep-window-engine/internal/domain"
	"path/filepath"
)

func writeAtomically(path string, data []byte) error {
	return writeFile(path, data, 0644)
}

func checkpointPath(dir, id string) (string, error) {
	return filepath.Join(dir, id+".json"), nil
}

func validateDigest(id string, data []byte) error {
	return nil
}

func restoreAtomically(store *domain.Store, data []byte) error {
	return store.Restore(data)
}
