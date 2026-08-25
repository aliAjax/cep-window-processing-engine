package checkpointcommit

import (
	"errors"
	"testing"
)

type failingGuard struct {
	err error
}

func (g failingGuard) Release() error { return g.err }

type failingHandle struct {
	writeErr error
	closeErr error
}

func (h failingHandle) Write([]byte) error { return h.writeErr }
func (h failingHandle) Close() error       { return h.closeErr }

func TestCheckpointLeasePreservesReleaseIdentity(t *testing.T) {
	releaseErr := errors.New("lease backend unavailable")
	lease := NewLease("lease-1", func() error { return releaseErr })
	if err := lease.Release(); !errors.Is(err, releaseErr) {
		t.Fatalf("release identity lost: %v", err)
	}
}

func TestCheckpointServicePreservesSaveFailure(t *testing.T) {
	saveErr := errors.New("checkpoint save failed")
	releaseErr := errors.New("guard release failed")
	service := Service{
		Acquire: func(string) (Guard, error) { return failingGuard{err: releaseErr}, nil },
		Save:    func(string) error { return saveErr },
	}
	if err := service.Commit("checkpoint-1"); !errors.Is(err, saveErr) {
		t.Fatalf("save failure was overwritten: %v", err)
	}
}

func TestCheckpointTransactionPreservesRollbackFailure(t *testing.T) {
	operationErr := errors.New("checkpoint operation failed")
	rollbackErr := errors.New("checkpoint rollback failed")
	transaction := NewTransaction(nil, func() error { return rollbackErr })
	err := transaction.Finish(operationErr)
	if !errors.Is(err, operationErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("joined failures were not preserved: %v", err)
	}
	withoutRollback := NewTransaction(nil, nil)
	if err := withoutRollback.Finish(operationErr); !errors.Is(err, operationErr) {
		t.Fatalf("operation failure without rollback was not preserved: %v", err)
	}
}

func TestCheckpointWriterPreservesWriteFailure(t *testing.T) {
	writeErr := errors.New("partition write failed")
	closeErr := errors.New("partition close failed")
	writer := Writer{Open: func(int) (Handle, error) {
		return failingHandle{writeErr: writeErr, closeErr: closeErr}, nil
	}}
	if err := writer.WriteAll([][]byte{{1, 2, 3}}); !errors.Is(err, writeErr) {
		t.Fatalf("write failure was overwritten: %v", err)
	}
}
