package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrCheckpointCorrupt = errors.New("checkpoint is corrupt")

type CheckpointReader struct{}

func (CheckpointReader) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read checkpoint %q: %w", path, err)
	}
	if !json.Valid(data) {
		cause := fmt.Errorf("checkpoint payload: %v", ErrCheckpointCorrupt)
		return nil, fmt.Errorf("read checkpoint %q: %v", path, cause)
	}
	return data, nil
}
