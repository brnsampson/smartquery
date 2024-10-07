package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/brnsampson/optional"
)

type MatchType int

const (
	// Matching operations defined. These are currently implemented for individual values, but could be extended to slices
	// of values as noted:
	MatchAlways MatchType = iota // This ALWAYS matches. It is true if S = S ∪ S which is always true.
	MatchNone                    // True if ⦰ = S
	MatchAny                     // True if ⦰ != S
	MatchSome                    // True if ⦰ != S1 ∩ S2. TODO: implement this! Until we support slices in queries this is the same as Exact though...
	MatchExact                   // True if ⦰ = S1 𝚫 S2
	MatchLike                    // Only valid for strings: perform
)

type Query[T comparable] interface {
	Matches(T) (bool, error)
	MatchesOption(optional.Optional[T]) (bool, error)
}

func Always[T comparable]() FieldQuery[T] {
	tmp := optional.None[T]()
	return FieldQuery[T]{MatchAlways, &tmp}
}

func None[T comparable](match T) FieldQuery[T] {
	tmp := optional.NewOption[T](match)
	return FieldQuery[T]{MatchNone, &tmp}
}

func Any[T comparable](match T) FieldQuery[T] {
	tmp := optional.NewOption[T](match)
	return FieldQuery[T]{MatchAny, &tmp}
}

func Exact[T comparable](match T) FieldQuery[T] {
	tmp := optional.NewOption[T](match)
	return FieldQuery[T]{MatchExact, &tmp}
}

func Like[T comparable](match T) FieldQuery[T] {
	tmp := optional.NewOption[T](match)
	return FieldQuery[T]{MatchLike, &tmp}
}

func AlwaysString() StringQuery {
	tmp := optional.None[string]()
	return StringQuery{MatchAlways, &tmp}
}

func NoneString(match string) StringQuery {
	tmp := optional.NewOption(match)
	return StringQuery{MatchNone, &tmp}
}

func AnyString(match string) StringQuery {
	tmp := optional.NewOption(match)
	return StringQuery{MatchAny, &tmp}
}

func ExactString(match string) StringQuery {
	tmp := optional.NewOption(match)
	return StringQuery{MatchExact, &tmp}
}

func LikeString(match string) StringQuery {
	tmp := optional.NewOption(match)
	return StringQuery{MatchLike, &tmp}
}

type FieldQuery[T comparable] struct {
	criteria MatchType
	value    optional.Optional[T]
}

func NewQuery[T comparable, O optional.Optional[T]](matchType MatchType, value O) FieldQuery[T] {
	return FieldQuery[T]{matchType, value}
}

func (q FieldQuery[T]) AsRef() *FieldQuery[T] {
	return &q
}

func (q *FieldQuery[T]) Matches(value T) (bool, error) {
	none := q.value.IsNone()
	tmp := q.value.Clone()
	var val T
	if !none {
		val = tmp.UnsafeUnwrap()
	}

	c := q.criteria
	if c == MatchAlways {
		return true, nil
	} else if c == MatchNone {
		return false, nil
	} else if c == MatchAny {
		// Non-option value passed, so it is always some value
		return true, nil
	} else if c == MatchExact {
		if none {
			// Just here so we never try to compare val if it is not initialized
			return false, nil
		} else if val == value {
			return true, nil
		} else {
			return false, nil
		}
	} else if c == MatchLike {
		// Not supported!
		return false, fmt.Errorf("QueryError: cannot perform MatchLike matches on generic type. Use StringQuery instead.")
	}
	return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
}

func (q *FieldQuery[T]) MatchesOption(value optional.Optional[T]) (bool, error) {
	none := q.value.IsNone()
	tmp := q.value.Clone()
	var val T
	if !none {
		val = tmp.UnsafeUnwrap()
	}

	c := q.criteria
	otherMatchNone := value.IsNone()
	otherMatchAny := !otherMatchNone

	if c == MatchAlways {
		return true, nil
	} else if (c == MatchNone && otherMatchNone) || (c == MatchAny && otherMatchAny) {
		return true, nil
	} else if (c == MatchNone && otherMatchAny) || (c == MatchAny && otherMatchNone) {
		return false, nil
	} else if c == MatchExact {
		other, _ := value.Clone().Unwrap()
		if none && otherMatchAny {
			// Just here so we never try to compare val if it is not initialized
			return false, nil
		} else if none && otherMatchNone {
			return true, nil
		} else if otherMatchNone {
			// q.value is SOME but the other value is NONE
			return false, nil
		} else if val == other {
			return true, nil
		} else {
			return false, nil
		}
	} else if c == MatchLike {
		// Not supported!
		return false, fmt.Errorf("QueryError: cannot perform MatchLike matches on generic type. Use StringQuery instead.")
	}
	return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
}

type StringQuery struct {
	criteria MatchType
	value    optional.Optional[string]
}

func NewStringQuery(matchType MatchType, value optional.Optional[string]) StringQuery {
	return StringQuery{matchType, value}
}

func (q StringQuery) AsRef() *StringQuery {
	return &q
}

func (q *StringQuery) Matches(value string) (bool, error) {
	c := q.criteria
	tmp := q.value.Clone()
	test, err := tmp.Unwrap()
	if err != nil {
		// q.value is MatchNone!
		if c == MatchAlways {
			return true, nil
		} else if c == MatchNone {
			return false, nil
		} else if c == MatchAny {
			// value is a non-option type, so it is always MatchAny
			return true, nil
		} else if c == MatchExact {
			// value is a non-option type, so it is always MatchAny
			return false, nil
		} else if c == MatchLike {
			// Not sure why this would come up, but I guess a MatchLike match of MatchNone matches nothing?
			return false, nil
		} else {
			return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
		}
	}

	// The case of q.value being MatchNone is handled above
	if c == MatchAlways {
		return true, nil
	} else if c == MatchNone {
		return false, nil
	} else if c == MatchAny {
		// Non-option value passed, so it is always some value
		return true, nil
	} else if c == MatchExact {
		if test == value {
			// SAFETY: we guarenteed that q.value is not MatchNone above
			return true, nil
		} else {
			return false, nil
		}
	} else if c == MatchLike {
		// MatchLike supports two wildcards, % for multiple characters and _ for a single char.
		// We support this by converting those to the regexp equivilants (.* and . respectively)
		tmp := strings.ReplaceAll(test, "%", ".*")
		tmp = strings.ReplaceAll(tmp, "_", ".")
		matches, err := regexp.MatchString(tmp, value)
		if err != nil {
			return false, err
		}
		return matches, nil
	}
	return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
}

func (q *StringQuery) MatchesOption(value optional.Optional[string]) (bool, error) {
	c := q.criteria
	tmp := q.value.Clone()
	test, err := tmp.Unwrap()

	otherMatchNone := value.IsNone()
	if err != nil {
		// q.value is MatchNone!
		if c == MatchAlways {
			return true, nil
		} else if c == MatchNone {
			return otherMatchNone, nil
		} else if c == MatchAny {
			return !otherMatchNone, nil
		} else if c == MatchExact {
			return otherMatchNone, nil
		} else if c == MatchLike {
			// Not sure why this would come up, but I guess a MatchLike match of MatchNone matches against MatchNone and nothing else?
			return otherMatchNone, nil
		} else {
			return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
		}
	}

	// The case of q.value being MatchNone is handled above
	other, err := value.Clone().Unwrap()
	if err != nil {
		// The case of value is MatchNone and q.value is MatchAny
		if c == MatchAlways {
			return true, nil
		} else if c == MatchNone {
			return true, nil
		} else if c == MatchAny {
			return false, nil
		} else if c == MatchExact {
			return false, nil
		} else if c == MatchLike {
			// MatchNone has no content that could possibly match
			return false, nil
		} else {
			return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
		}
	}

	// The case of value and q.value are both MatchAny
	if c == MatchAlways {
		return true, nil
	} else if c == MatchNone {
		return false, nil
	} else if c == MatchAny {
		return true, nil
	} else if c == MatchExact {
		if test == other {
			// SAFETY: we guarenteed that q.value and value are both MatchAny above
			return true, nil
		} else {
			return false, nil
		}
	} else if c == MatchLike {
		// MatchLike supports two wildcards, % for multiple characters and _ for a single char.
		// We support this by converting those to the regexp equivilants (.* and . respectively)
		tmp := strings.ReplaceAll(test, "%", ".*")
		tmp = strings.ReplaceAll(tmp, "_", ".")
		matches, err := regexp.MatchString(tmp, other)
		if err != nil {
			return false, err
		}
		return matches, nil
	}
	return false, fmt.Errorf("QueryError: unsupported matching strategy: %d", c)
}
