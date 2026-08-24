package quota

import "fmt"

type QuotaApplier interface {
	Apply(Reservation) error
}

type QuotaClassifier interface {
	Decide(error) Decision
}

type QuotaHandlerError struct {
	Cause error
}

func (e QuotaHandlerError) Error() string {
	return fmt.Sprintf("submit quota reservation: %v", e.Cause)
}

type Response struct {
	Status   int
	Decision Decision
	Err      error
}

type Handler struct {
	service QuotaApplier
	policy  QuotaClassifier
}

func NewHandler(service QuotaApplier, policy QuotaClassifier) *Handler {
	return &Handler{service: service, policy: policy}
}

func (h *Handler) Submit(reservation Reservation) Response {
	if err := h.service.Apply(reservation); err != nil {
		exposed := QuotaHandlerError{Cause: err}
		decision := h.policy.Decide(exposed)
		status := 503
		if decision == DecisionReject {
			status = 429
		}
		return Response{Status: status, Decision: decision, Err: exposed}
	}
	return Response{Status: 201}
}
