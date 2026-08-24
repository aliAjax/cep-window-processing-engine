package quota

import "fmt"

type Reserver interface {
	Reserve(Reservation) error
}

type QuotaServiceError struct {
	TenantID string
	Cause    error
}

func (e QuotaServiceError) Error() string {
	return fmt.Sprintf("apply tenant quota %q: %v", e.TenantID, e.Cause)
}

type Service struct {
	repository Reserver
}

func NewService(repository Reserver) *Service {
	return &Service{repository: repository}
}

func (s *Service) Apply(reservation Reservation) error {
	if err := s.repository.Reserve(reservation); err != nil {
		return QuotaServiceError{TenantID: reservation.TenantID, Cause: err}
	}
	return nil
}
