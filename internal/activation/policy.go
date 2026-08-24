package activation

import (
	"errors"
	"fmt"
)

type Decision string

const (
	DecisionReject Decision = "reject"
	DecisionRetry  Decision = "retry"
)

type PolicyError struct {
	Cause error
}

func (e PolicyError) Error() string {
	return fmt.Sprintf("classify activation failure: %v", e.Cause)
}

type Policy struct{}

func (Policy) Decide(err error) Decision {
	classified := PolicyError{Cause: err}
	if errors.Is(classified, ErrInvalidPlan) {
		return DecisionReject
	}
	return DecisionRetry
}
