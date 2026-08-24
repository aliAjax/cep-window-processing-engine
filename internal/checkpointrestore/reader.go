package checkpointrestore

import (
	"errors"
	"fmt"
)

var (
	ErrCheckpointMissing = errors.New("checkpoint missing")
	ErrCheckpointCorrupt = errors.New("checkpoint corrupt")
)

type ReadError struct {
	CheckpointID string
	Err          error
}

func (e *ReadError) Error() string {
	return fmt.Sprintf("read checkpoint %q: %v", e.CheckpointID, e.Err)
}

func (e *ReadError) Unwrap() error {
	return e.Err
}

type Reader struct {
	Blobs map[string][]byte
}

func (r Reader) Read(id string) ([]byte, error) {
	raw, ok := r.Blobs[id]
	if !ok {
		return nil, &ReadError{CheckpointID: id, Err: ErrCheckpointMissing}
	}
	if len(raw) == 0 {
		return nil, &ReadError{CheckpointID: id, Err: ErrCheckpointCorrupt}
	}
	return append([]byte(nil), raw...), nil
}
