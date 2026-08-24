package checkpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/cep-window-engine/internal/domain"
	"os"
	"path/filepath"
)

type Manager struct{ Dir string }

func New(dir string) *Manager { return &Manager{Dir: dir} }
func (m *Manager) Save(s *domain.Store) (string, error) {
	if err := os.MkdirAll(m.Dir, 0755); err != nil {
		return "", err
	}
	b, err := s.Snapshot()
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	id := hex.EncodeToString(h[:])
	if err = os.WriteFile(filepath.Join(m.Dir, id+".json"), b, 0644); err != nil {
		return "", err
	}
	return id, nil
}
func (m *Manager) Restore(s *domain.Store, id string) error {
	b, e := os.ReadFile(filepath.Join(m.Dir, id+".json"))
	if e != nil {
		return fmt.Errorf("checkpoint: %w", e)
	}
	return s.Restore(b)
}
