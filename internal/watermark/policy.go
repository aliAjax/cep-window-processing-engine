package watermark

import "time"

type Policy struct {
	MaximumStep time.Duration
}

func (p Policy) Accept(current, next time.Time) bool {
	return current.IsZero() || !next.Equal(current)
}
