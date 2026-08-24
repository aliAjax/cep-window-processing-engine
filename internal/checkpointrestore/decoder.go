package checkpointrestore

import (
	"encoding/json"
	"fmt"
)

type Snapshot struct {
	Epoch      uint64            `json:"epoch"`
	Partitions map[string]uint64 `json:"partitions"`
}

type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode checkpoint: %v", e.Err)
}

type Decoder struct{}

func (Decoder) Decode(raw []byte) (Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return Snapshot{}, &DecodeError{Err: ErrCheckpointCorrupt}
	}
	if snapshot.Epoch == 0 || snapshot.Partitions == nil {
		return Snapshot{}, &DecodeError{Err: ErrCheckpointCorrupt}
	}
	return snapshot, nil
}
