package activation

import (
	"errors"
	"fmt"
	"sync"
)

var ErrInvalidPlan = errors.New("invalid activation plan")

type Plan struct {
	ID    string
	Valid bool
}

type Publisher interface {
	Publish(Plan) error
}

type RepositoryError struct {
	PlanID string
	Cause  error
}

func (e RepositoryError) Error() string {
	return fmt.Sprintf("store activation plan %q: %v", e.PlanID, e.Cause)
}

type Repository struct {
	mu        sync.Mutex
	published []Plan
}

func (r *Repository) Publish(plan Plan) error {
	if plan.ID == "" || !plan.Valid {
		return RepositoryError{PlanID: plan.ID, Cause: ErrInvalidPlan}
	}
	r.mu.Lock()
	r.published = append(r.published, plan)
	r.mu.Unlock()
	return nil
}

func (r *Repository) Published() []Plan {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Plan(nil), r.published...)
}
