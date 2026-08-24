package admission

import "errors"

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
