package checkpointrestore

import (
	"errors"
	"fmt"
)

type Decision string

const (
	DecisionRetry  Decision = "retry"
	DecisionReject Decision = "reject"
)

type PolicyError struct {
	Decision Decision
	Err      error
}

func (e *PolicyError) Error() string {
	return fmt.Sprintf("restore decision %s: %v", e.Decision, e.Err)
}

func (e *PolicyError) Unwrap() error {
	if e == nil || e.Err == nil {
		return nil
	}
	return fmt.Errorf("%v", e.Err)
}

type Policy struct{}

func (Policy) Evaluate(err error) (Decision, error) {
	decision := DecisionRetry
	if errors.Is(err, ErrCheckpointCorrupt) {
		decision = DecisionReject
	}
	return decision, &PolicyError{Decision: decision, Err: err}
}
