package partition

import "encoding/json"

type Snapshot struct {
	Assignments map[int]string `json:"assignments"`
	Generation  uint64         `json:"generation"`
}

func Decode(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(data) > 0 {
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return Snapshot{}, err
		}
	}
	if snapshot.Assignments == nil {
		snapshot.Assignments = make(map[int]string)
	}
	return snapshot, nil
}
