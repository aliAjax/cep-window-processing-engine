package quota

import (
	"errors"
	"fmt"
)

type Decision string

const (
	DecisionReject Decision = "reject"
	DecisionRetry  Decision = "retry"
)

type QuotaPolicyError struct {
	Cause error
}

func (e QuotaPolicyError) Error() string {
	return fmt.Sprintf("classify quota reservation: %v", e.Cause)
}

type Policy struct{}

func (Policy) Decide(err error) Decision {
	classified := QuotaPolicyError{Cause: err}
	if errors.Is(classified, ErrQuotaExceeded) {
		return DecisionReject
	}
	return DecisionRetry
}
