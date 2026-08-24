package activation

import "fmt"

type Activator interface {
	Activate(Plan) error
}

type Classifier interface {
	Decide(error) Decision
}

type HandlerError struct {
	Cause error
}

func (e HandlerError) Error() string {
	return fmt.Sprintf("publish activation: %v", e.Cause)
}

type Response struct {
	Status   int
	Decision Decision
	Err      error
}

type Handler struct {
	service Activator
	policy  Classifier
}

func NewHandler(service Activator, policy Classifier) *Handler {
	return &Handler{service: service, policy: policy}
}

func (h *Handler) Publish(plan Plan) Response {
	if err := h.service.Activate(plan); err != nil {
		exposed := HandlerError{Cause: err}
		decision := h.policy.Decide(exposed)
		status := 503
		if decision == DecisionReject {
			status = 422
		}
		return Response{Status: status, Decision: decision, Err: exposed}
	}
	return Response{Status: 201}
}
