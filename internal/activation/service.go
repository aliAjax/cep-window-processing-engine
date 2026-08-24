package activation

import "fmt"

type ServiceError struct {
	PlanID string
	Cause  error
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("activate plan %q: %v", e.PlanID, e.Cause)
}

type Service struct {
	publisher Publisher
}

func NewService(publisher Publisher) *Service {
	return &Service{publisher: publisher}
}

func (s *Service) Activate(plan Plan) error {
	if err := s.publisher.Publish(plan); err != nil {
		return ServiceError{PlanID: plan.ID, Cause: err}
	}
	return nil
}
