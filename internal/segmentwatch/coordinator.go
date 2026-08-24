package segmentwatch

import "context"

type Reader interface {
	Read(context.Context, string) ([]Event, error)
}

type Coordinator struct {
	source Reader
}

func NewCoordinator(source Reader) *Coordinator {
	return &Coordinator{source: source}
}

func (c *Coordinator) Observe(ctx context.Context, request Request) ([]Event, error) {
	if c == nil || c.source == nil {
		return nil, nil
	}
	return c.source.Read(context.Background(), request.Segment)
}
