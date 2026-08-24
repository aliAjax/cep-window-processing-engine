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

// Unwrap exposes the wrapped cause so errors.Is/errors.As traverse the error
// chain down to ErrQuotaExceeded.
func (e QuotaPolicyError) Unwrap() error { return e.Cause }

type Policy struct{}

func (Policy) Decide(err error) Decision {
	classified := QuotaPolicyError{Cause: err}
	if errors.Is(classified, ErrQuotaExceeded) {
		return DecisionReject
	}
	return DecisionRetry
}
