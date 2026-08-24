package delivery

import (
	"context"
	"fmt"

	"github.com/example/cep-window-engine/internal/domain"
)

type Sender interface {
	Send(context.Context, domain.Match) error
}

type Service struct {
	Sender Sender
}

func (s Service) Deliver(ctx context.Context, match domain.Match, lease *Lease) error {
	if lease == nil {
		return fmt.Errorf("delivery lease is required")
	}
	if s.Sender == nil {
		return fmt.Errorf("delivery sender is unavailable")
	}
	if err := s.Sender.Send(ctx, match); err != nil {
		return fmt.Errorf("deliver match %s: %w", match.ID, err)
	}
	if err := lease.Release(); err != nil {
		return fmt.Errorf("release delivery lease: %w", err)
	}
	return nil
}
