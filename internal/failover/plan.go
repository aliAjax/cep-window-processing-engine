package failover

import "sync"

type Step struct {
	ID string
}

type Plan struct {
	mu     sync.Mutex
	state  State
	steps  []Step
	active chan struct{}
	once   sync.Once
}

func NewPlan(steps ...Step) *Plan {
	return &Plan{state: StateStandby, steps: append([]Step(nil), steps...), active: make(chan struct{})}
}

func (p *Plan) Begin() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.state = StatePromoting
}

func (p *Plan) Next() (Step, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.steps) == 0 {
		return Step{}, false
	}
	step := p.steps[0]
	p.steps = p.steps[1:]
	return step, true
}

func (p *Plan) State() State {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state
}

func (p *Plan) markActive() {
	p.mu.Lock()
	p.state = StateActive
	p.mu.Unlock()
	p.once.Do(func() { close(p.active) })
}

func (p *Plan) activeSignal() <-chan struct{} {
	return p.active
}
