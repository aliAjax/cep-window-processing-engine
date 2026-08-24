package compiler

import (
	"crypto/sha256"
	"fmt"
	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/pattern/ast"
	"strings"
	"time"
)

type Compiler struct{ MaxComplexity int }

func (c Compiler) Compile(n ast.Node) (*domain.CompiledPlan, error) {
	if c.MaxComplexity == 0 {
		c.MaxComplexity = 100
	}
	steps := []domain.PlanStep{}
	complexity := 0
	var walk func(ast.Node) error
	walk = func(x ast.Node) error {
		complexity++
		if complexity > c.MaxComplexity {
			return fmt.Errorf("pattern complexity exceeds quota")
		}
		switch v := x.(type) {
		case ast.Predicate:
			steps = append(steps, domain.PlanStep{Alias: v.EventType, EventType: v.EventType, Predicate: v.Expression})
		case ast.Sequence:
			for _, i := range v.Items {
				if e := walk(i); e != nil {
					return e
				}
			}
		case ast.Not:
			if e := walk(v.Item); e != nil {
				return e
			}
			steps[len(steps)-1].Negated = true
		case ast.Within:
			d, e := time.ParseDuration(v.Duration)
			if e != nil {
				return e
			}
			if e = walk(v.Item); e != nil {
				return e
			}
			steps[len(steps)-1].Within = d
		case ast.Repeat:
			if e := walk(v.Item); e != nil {
				return e
			}
			steps[len(steps)-1].Repeat = v.Max
		case ast.Concurrent, ast.FollowedBy, ast.Join:
			return fmt.Errorf("operator not supported in compact parser")
		}
		return nil
	}
	if err := walk(n); err != nil {
		return nil, err
	}
	raw := strings.Builder{}
	for _, s := range steps {
		raw.WriteString(s.EventType + s.Predicate)
	}
	h := sha256.Sum256([]byte(raw.String()))
	return &domain.CompiledPlan{Steps: steps, Complexity: complexity, Hash: fmt.Sprintf("%x", h[:])}, nil
}
