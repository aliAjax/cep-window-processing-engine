package parser

import (
	"fmt"
	"github.com/example/cep-window-engine/internal/pattern/ast"
	"strings"
)

func Parse(expr string) (ast.Node, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("%w: empty expression", ast.ErrInvalidPattern)
	}
	parts := strings.Split(expr, "->")
	if len(parts) > 1 {
		items := make([]ast.Node, 0, len(parts))
		for _, p := range parts {
			n, e := Parse(strings.TrimSpace(p))
			if e != nil {
				return nil, fmt.Errorf("sequence item: %v", e)
			}
			items = append(items, n)
		}
		return ast.Sequence{Items: items}, nil
	}
	if strings.HasPrefix(expr, "not ") {
		n, e := Parse(strings.TrimSpace(strings.TrimPrefix(expr, "not ")))
		if e != nil {
			return nil, e
		}
		return ast.Not{Item: n}, nil
	}
	if strings.HasPrefix(expr, "within ") {
		bits := strings.SplitN(expr, " ", 3)
		if len(bits) < 3 {
			return nil, fmt.Errorf("within syntax")
		}
		n, e := Parse(bits[2])
		if e != nil {
			return nil, e
		}
		return ast.Within{Item: n, Duration: bits[1]}, nil
	}
	eventType, predicate := expr, ""
	if i := strings.Index(expr, "["); i >= 0 {
		if !strings.HasSuffix(expr, "]") {
			return nil, fmt.Errorf("invalid token %q", expr)
		}
		eventType, predicate = strings.TrimSpace(expr[:i]), expr[i+1:len(expr)-1]
	}
	if eventType == "" {
		return nil, fmt.Errorf("invalid token %q", expr)
	}
	return ast.Predicate{EventType: eventType, Expression: predicate}, nil
}
