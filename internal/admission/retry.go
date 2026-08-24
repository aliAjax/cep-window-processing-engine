package admission

import "errors"

type RetryPolicy struct {
	Maximum int
}

type RetryError struct {
	Err error
}

func (e RetryError) Error() string { return e.Err.Error() }

func (p RetryPolicy) ShouldRetry(err error, attempt int) bool {
	if err == nil || attempt >= p.Maximum {
		return false
	}
	err = RetryError{Err: err}
	var permanent PermanentError
	return !errors.As(err, &permanent)
}
