package quota

import (
	"errors"
	"fmt"
	"testing"
)

type reserverFunc func(Reservation) error

func (f reserverFunc) Reserve(reservation Reservation) error { return f(reservation) }

type quotaApplierFunc func(Reservation) error

func (f quotaApplierFunc) Apply(reservation Reservation) error { return f(reservation) }

type quotaClassifierFunc func(error) Decision

func (f quotaClassifierFunc) Decide(err error) Decision { return f(err) }

func TestQuotaRepositoryPreservesExceeded(t *testing.T) {
	repository := NewRepository(map[string]int{"tenant-a": 2})
	err := repository.Reserve(Reservation{TenantID: "tenant-a", Units: 3})
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("repository lost quota-exceeded sentinel: %v", err)
	}
}

func TestQuotaServicePreservesRepositoryError(t *testing.T) {
	service := NewService(reserverFunc(func(Reservation) error { return ErrQuotaExceeded }))
	err := service.Apply(Reservation{TenantID: "tenant-a", Units: 3})
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("service lost repository error chain: %v", err)
	}
}

func TestQuotaPolicyRejectsExceeded(t *testing.T) {
	if decision := (Policy{}).Decide(ErrQuotaExceeded); decision != DecisionReject {
		t.Fatalf("quota excess classified as %q, want reject", decision)
	}
}

func TestQuotaHandlerRejectsExceededWithoutRetry(t *testing.T) {
	service := quotaApplierFunc(func(Reservation) error { return ErrQuotaExceeded })
	policy := quotaClassifierFunc(func(err error) Decision {
		if errors.Is(err, ErrQuotaExceeded) {
			return DecisionReject
		}
		return DecisionRetry
	})
	response := NewHandler(service, policy).Submit(Reservation{TenantID: "tenant-a", Units: 3})
	if response.Status != 429 || response.Decision != DecisionReject {
		panic(fmt.Sprintf("quota excess mapped to status=%d decision=%s", response.Status, response.Decision))
	}
}
