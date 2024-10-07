package query

import (
	"github.com/brnsampson/optional"
)

type Matcher interface {
	Match() (bool, error)
}

type Match[T comparable] struct {
	query   Query[T]
	operand optional.Optional[T]
}

func NewValueMatch[T comparable](operand T, query Query[T]) Matcher {
	return &Match[T]{query, optional.NewOption(operand).AsRef()}
}

func NewMatch[T comparable](operand optional.Optional[T], query Query[T]) Matcher {
	return &Match[T]{query, operand}
}

func (m Match[T]) Match() (bool, error) {
	return m.query.MatchesOption(m.operand)
}

// Helper function for when the user is creating a custom query type. Create all the Match objects, dump them in a
// slice, then use MatchAll.
func MatchAll(matches []Matcher) (bool, error) {
	for _, m := range matches {
		matched, err := m.Match()
		if err != nil {
			return false, err
		}
		if matched == false {
			return false, nil
		}
	}
	return true, nil
}
