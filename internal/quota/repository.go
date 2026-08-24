package quota

import (
	"errors"
	"fmt"
	"sync"
)

var ErrQuotaExceeded = errors.New("quota reservation exceeds available budget")

type Reservation struct {
	TenantID string
	Units    int
}

type QuotaRepositoryError struct {
	TenantID string
	Cause    error
}

func (e QuotaRepositoryError) Error() string {
	return fmt.Sprintf("reserve quota for tenant %q: %v", e.TenantID, e.Cause)
}

type Repository struct {
	mu        sync.Mutex
	available map[string]int
}

func NewRepository(available map[string]int) *Repository {
	return &Repository{available: available}
}

func (r *Repository) Reserve(reservation Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reservation.Units <= 0 || reservation.Units > r.available[reservation.TenantID] {
		return QuotaRepositoryError{TenantID: reservation.TenantID, Cause: ErrQuotaExceeded}
	}
	r.available[reservation.TenantID] -= reservation.Units
	return nil
}
