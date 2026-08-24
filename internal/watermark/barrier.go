package watermark

import (
	"time"
)

type Barrier struct {
	current time.Time
}

func (b *Barrier) Publish(candidate time.Time) time.Time {
	if candidate.After(b.current) {
		b.current = candidate
	}
	return b.current
}

func (b *Barrier) Current() time.Time {
	return b.current
}
