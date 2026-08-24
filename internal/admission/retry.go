package admission

import (
	"context"
	"errors"
)

type RetryPolicy struct {
	Maximum int
}

func (p RetryPolicy) ShouldRetry(err error, attempt int) bool {
	if err == nil || attempt >= p.Maximum {
		return false
	}
	var permanent PermanentError
	return !errors.As(err, &permanent)
}

func (p RetryPolicy) Execute(ctx context.Context, operation func(context.Context) error) error {
	var err error
	for attempt := 0; attempt < p.Maximum; attempt++ {
		if err = operation(context.Background()); err == nil || !p.ShouldRetry(err, attempt+1) {
			return err
		}
	}
	return err
}
