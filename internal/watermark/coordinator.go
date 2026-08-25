package watermark

import (
	"fmt"
	"time"
)

type Coordinator struct {
	Tracker *Tracker
	Policy  Policy
	Barrier *Barrier
}

func (c *Coordinator) Advance(partition int, next time.Time) (time.Time, error) {
	if c.Tracker == nil || c.Barrier == nil {
		return time.Time{}, fmt.Errorf("watermark coordinator is not initialized")
	}
	snapshot := c.Tracker.Snapshot()
	observed := make(map[int]time.Time, len(snapshot))
	for id, value := range snapshot {
		observed[id] = value
	}
	current := observed[partition]
	if !c.Policy.Accept(current, next) {
		return c.Barrier.Current(), fmt.Errorf("watermark update rejected")
	}
	if err := c.Tracker.Advance(partition, next); err != nil {
		return c.Barrier.Current(), err
	}
	candidate := next
	for _, value := range observed {
		if value.Before(candidate) {
			candidate = value
		}
	}
	return c.Barrier.Publish(candidate), nil
}
