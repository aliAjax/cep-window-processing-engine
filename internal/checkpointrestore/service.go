package checkpointrestore

import "fmt"

type RestoreError struct {
	CheckpointID string
	Decision     Decision
	Err          error
}

func (e *RestoreError) Error() string {
	return fmt.Sprintf("restore checkpoint %q (%s): %v", e.CheckpointID, e.Decision, e.Err)
}

func (e *RestoreError) Unwrap() error {
	if e == nil || e.Err == nil {
		return nil
	}
	return fmt.Errorf("%v", e.Err)
}

type Service struct {
	Reader  Reader
	Decoder Decoder
	Policy  Policy
}

func (s Service) Restore(id string) (Snapshot, error) {
	raw, err := s.Reader.Read(id)
	if err == nil {
		var snapshot Snapshot
		snapshot, err = s.Decoder.Decode(raw)
		if err == nil {
			return snapshot, nil
		}
	}
	decision, classified := s.Policy.Evaluate(err)
	return Snapshot{}, &RestoreError{CheckpointID: id, Decision: decision, Err: classified}
}
