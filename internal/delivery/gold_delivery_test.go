package delivery

import (
	"context"
	"errors"
	"testing"

	"github.com/example/cep-window-engine/internal/domain"
)

type senderFunc func(context.Context, domain.Match) error

func (f senderFunc) Send(ctx context.Context, match domain.Match) error {
	return f(ctx, match)
}

type transactionStub struct {
	commitErr   error
	rollbackErr error
}

func (t transactionStub) Commit() error   { return t.commitErr }
func (t transactionStub) Rollback() error { return t.rollbackErr }

func TestBatchReleasesEachHandle(t *testing.T) {
	var previous *Lease
	acquire := func() (*Lease, error) {
		if previous != nil && !previous.Released() {
			return nil, errors.New("delivery leases exhausted")
		}
		previous = &Lease{}
		return previous, nil
	}
	batch := Batch{Service: Service{Sender: senderFunc(func(context.Context, domain.Match) error { return nil })}}
	matches := []domain.Match{{ID: "m-1"}, {ID: "m-2"}, {ID: "m-3"}}
	if err := batch.SendAll(context.Background(), matches, acquire); err != nil {
		t.Fatalf("send batch: %v", err)
	}
}

func TestTransactionPreservesBusinessError(t *testing.T) {
	businessErr := errors.New("sink rejected match")
	rollbackErr := errors.New("rollback failed")
	err := Finish(transactionStub{rollbackErr: rollbackErr}, businessErr)
	if !errors.Is(err, businessErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("finish lost error chain: %v", err)
	}
}

func TestServiceRollsBackFailedDelivery(t *testing.T) {
	lease := &Lease{}
	sinkErr := errors.New("sink unavailable")
	service := Service{Sender: senderFunc(func(context.Context, domain.Match) error { return sinkErr })}
	if err := service.Deliver(context.Background(), domain.Match{ID: "m-fail"}, lease); !errors.Is(err, sinkErr) {
		t.Fatalf("expected sink error, got %v", err)
	}
	if !lease.Released() {
		t.Fatal("failed delivery retained its lease")
	}
}

func TestLeaseReleasedOnEarlyReturn(t *testing.T) {
	lease := &Lease{}
	if err := lease.Release(); err != nil {
		t.Fatalf("release lease: %v", err)
	}
	if !lease.Released() {
		t.Fatal("release did not update lease state")
	}
}
