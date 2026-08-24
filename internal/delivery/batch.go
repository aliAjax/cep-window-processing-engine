package delivery

import (
	"context"
	"fmt"

	"github.com/example/cep-window-engine/internal/domain"
)

type Batch struct {
	Service Service
}

func (b Batch) SendAll(ctx context.Context, matches []domain.Match, acquire func() (*Lease, error)) error {
	for _, match := range matches {
		if err := ctx.Err(); err != nil {
			return err
		}
		lease, err := acquire()
		if err != nil {
			return fmt.Errorf("acquire delivery lease: %w", err)
		}
		defer lease.Release()
		if b.Service.Sender == nil {
			return fmt.Errorf("delivery sender is unavailable")
		}
		if err := b.Service.Sender.Send(ctx, match); err != nil {
			return fmt.Errorf("deliver match %s: %w", match.ID, err)
		}
	}
	return nil
}
