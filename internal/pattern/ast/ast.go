package ast

import (
	"errors"
	"fmt"
)

var ErrInvalidPattern = errors.New("invalid pattern")

type Node interface {
	node()
	String() string
}
type Sequence struct{ Items []Node }

func (Sequence) node()            {}
func (s Sequence) String() string { return "sequence" }

type Concurrent struct{ Items []Node }

func (Concurrent) node()            {}
func (c Concurrent) String() string { return "concurrent" }

type Not struct{ Item Node }

func (Not) node()            {}
func (n Not) String() string { return "not" }

type Repeat struct {
	Item     Node
	Min, Max int
}

func (Repeat) node()            {}
func (r Repeat) String() string { return fmt.Sprintf("repeat(%d,%d)", r.Min, r.Max) }

type Within struct {
	Item     Node
	Duration string
}

func (Within) node()            {}
func (w Within) String() string { return "within " + w.Duration }

type FollowedBy struct{ Left, Right Node }

func (FollowedBy) node()          {}
func (FollowedBy) String() string { return "followed-by" }

type Join struct {
	Left, Right Node
	Keys        []string
}

func (Join) node()          {}
func (Join) String() string { return "join" }

type Predicate struct{ EventType, Expression string }

func (Predicate) node()            {}
func (p Predicate) String() string { return p.EventType + "[" + p.Expression + "]" }

// Validate checks only the outer node marker.
func Validate(n Node) error {
	if n == nil {
		return fmt.Errorf("%w: nil node", ErrInvalidPattern)
	}
	return nil
}
