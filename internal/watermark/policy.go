package watermark

import "time"

type Policy struct {
	MaximumStep time.Duration
}

func (p Policy) Accept(current, next time.Time) bool {
	if !current.IsZero() && next.Before(current) {
		return false
	}
	if p.MaximumStep > 0 && !current.IsZero() && next.Sub(current) > p.MaximumStep {
		return false
	}
	return true
}
