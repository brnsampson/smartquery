package query_test

import (
	"testing"

	"github.com/brnsampson/optional"
	"github.com/brnsampson/smartquery"
	"gotest.tools/v3/assert"
)

type testStruct struct {
	Name    string
	Email   optional.Optional[string]
	Balance int
	Stars   optional.Optional[int]
}

type testStructQuery struct {
	Name    smartquery.Query[string]
	Email   smartquery.Query[string]
	Balance smartquery.FieldQuery[int]
	Stars   smartquery.FieldQuery[int]
}

func (q *testStructQuery) Matches(test testStruct) (bool, error) {
	name, err := q.Name.Matches(test.Name)
	if err != nil {
		return false, err
	}
	email, err := q.Email.MatchesOption(test.Email)
	if err != nil {
		return false, err
	}
	balance, err := q.Balance.Matches(test.Balance)
	if err != nil {
		return false, err
	}
	stars, err := q.Stars.MatchesOption(test.Stars)
	if err != nil {
		return false, err
	}
	if name && email && balance && stars {
		return true, nil
	} else {
		return false, nil
	}
}

func (q *testStructQuery) MatchesOption(test optional.Optional[testStruct]) (bool, error) {
	if test.IsNone() {
		return false, nil
	}

	return q.Matches(test.UnsafeUnwrap())
}

func TestInterfaceStructQuery(t *testing.T) {
	name := "Chester the Tester"
	s := testStruct{Name: name, Email: optional.NewOption("chester@testing.org").AsRef(), Balance: 42, Stars: optional.NewOption(7).AsRef()}
	q := testStructQuery{
		Name:    smartquery.ExactString(name).AsRef(),
		Email:   smartquery.AlwaysString().AsRef(),
		Balance: smartquery.Any(42),
		Stars:   smartquery.Always[int](),
	}

	var tmp smartquery.Query[testStruct] = &q
	m, err := tmp.Matches(s)
	assert.NilError(t, err, "Error while trying to perform query match on custom testStruct!")
	assert.Assert(t, m, "Custom query on testStruct did not match!")
}

func TestAlwaysFieldQuery(t *testing.T) {
	original := 47
	originalOption := optional.NewOption(original)
	none := optional.None[int]()

	originalOptionCopy := optional.NewOption(original)
	wrongOption := optional.NewOption(42)
	noneCopy := optional.None[int]()

	// we expect everyone to use smartquery.Always(), but since a user could technically create the query in any of
	// these states we should test them all. Note that Always() is the same as empty below.
	different := smartquery.NewQuery(smartquery.MatchAlways, &wrongOption)
	same := smartquery.NewQuery(smartquery.MatchAlways, &originalOptionCopy)
	empty := smartquery.NewQuery(smartquery.MatchAlways, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match something!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[same] query does not match something!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[same] query does not match option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[same] query does not match option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match something!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match option with None value!")
}

func TestAlwaysStringQuery(t *testing.T) {
	original := "original"
	originalOption := optional.NewOption(original)
	originalOptionCopy := optional.NewOption(original)
	changed := "changed"
	changedOption := optional.NewOption(changed)
	none := optional.None[string]()
	noneCopy := optional.None[string]()

	// we expect everyone to use smartquery.AlwaysString(), but since a user could technically create the query in any of
	// these states we should test them all. Note that AlwaysString() is the same as empty below.
	different := smartquery.NewStringQuery(smartquery.MatchAlways, &changedOption)
	same := smartquery.NewStringQuery(smartquery.MatchAlways, &originalOptionCopy)
	empty := smartquery.NewStringQuery(smartquery.MatchAlways, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match something!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[some] query does not match something!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[different] query does not match option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match something!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Always[None] query does not match option with None value!")
}

func TestNoneFieldQuery(t *testing.T) {
	original := 47
	changed := 42
	originalOption := optional.NewOption(original)
	none := optional.None[int]()
	noneCopy := optional.None[int]()

	different := smartquery.None(changed)
	same := smartquery.None(original)
	empty := smartquery.NewQuery(smartquery.MatchNone, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[different] query matches something with a value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[different] query matches an option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[different] query does not match an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[same] query matches something with a value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[same] query matches an option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[same] query does not match an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[None] query matches something with a value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[None] query matches an option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[None] query does not match an option with None value!")
}

func TestNoneStringQuery(t *testing.T) {
	original := "original"
	changed := "changed"
	originalOption := optional.NewOption(original)
	none := optional.None[string]()
	noneCopy := optional.None[string]()

	different := smartquery.NoneString(changed)
	same := smartquery.NoneString(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewStringQuery(smartquery.MatchNone, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[different] query matches something with a value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[different] query matches an option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[different] query does not match an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[same] query matches something with a value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[same] query matches an option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[same] query does not match an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[None] query matches something with a value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "None[None] query matches an option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "None[None] query does not match an option with None value!")
}

func TestAnyFieldQuery(t *testing.T) {
	original := 47
	changed := 42
	originalOption := optional.NewOption(original)
	none := optional.None[int]()
	noneCopy := optional.None[int]()

	different := smartquery.Any[int](changed)
	same := smartquery.Any[int](original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewQuery[int](smartquery.MatchAny, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[different] query does not match something with a value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[different] query does not match option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[different] query matches an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[same] query does not match something with a value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[same] query does not match option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[same] query matches an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[None] query does not match something with a value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[None] query does not match option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[None] query matches an option with None value!")
}

func TestAnyStringQuery(t *testing.T) {
	original := "original"
	changed := "changed"
	originalOption := optional.NewOption(original)
	none := optional.None[string]()
	noneCopy := optional.None[string]()

	different := smartquery.AnyString(changed)
	same := smartquery.AnyString(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewStringQuery(smartquery.MatchAny, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[different] query does not match something with a value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[different] query does not match option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[different] query matches an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[same] query does not match something with a value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[same] query does not match option with Any value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[same] query matches an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[None] query does not match something with a value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Any[None] query does not match option with Any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Any[None] query matches an option with None value!")
}

func TestExactFieldQuery(t *testing.T) {
	original := 47
	changed := 42
	originalOption := optional.NewOption(original)
	none := optional.None[int]()
	noneCopy := optional.None[int]()

	different := smartquery.Exact(changed)
	same := smartquery.Exact(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewQuery(smartquery.MatchExact, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched something with a different value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched a Any option with a different value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched a value with an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[same] query did not match something with the same value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[same] query did not match a Any option with the same value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[same] query matched a value with an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[None] query matched something other than an option with None value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[None] query matches a Any option!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[None] query did not match an option with None value!")
}

func TestExactStringQuery(t *testing.T) {
	original := "original"
	changed := "changed"
	originalOption := optional.NewOption(original)
	none := optional.None[string]()
	noneCopy := optional.None[string]()

	different := smartquery.ExactString(changed)
	same := smartquery.ExactString(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewStringQuery(smartquery.MatchExact, &noneCopy)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched something with a different value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched a Any option with a different value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[different] query matched a value with an option with None value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[same] query did not match something with the same value!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[same] query did not match a Any option with the same value!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[same] query matched a value with an option with None value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[None] query matched something other than an option with None value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Exact[None] query matches a Any option!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Exact[None] query did not match an option with None value!")
}

func TestLikeFieldQuery(t *testing.T) {
	original := 47
	changed := 42
	originalPrefix := 4
	originalOption := optional.NewOption(original)
	none := optional.None[int]()
	noneCopy := optional.None[int]()
	prefixOption := optional.NewOption(originalPrefix)

	different := smartquery.Like(changed)
	same := smartquery.Like(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewQuery(smartquery.MatchLike, &noneCopy)
	prefix := smartquery.Like(originalPrefix)

	matches, err := different.Matches(original)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[different] query should not be supported for FieldQuery!")

	matches, err = different.MatchesOption(&originalOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[different] query should not be supported for FieldQuery!")

	matches, err = different.MatchesOption(&none)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[different] query should not be supported for FieldQuery!")

	matches, err = different.Matches(originalPrefix)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[different] query should not be supported for FieldQuery!")

	matches, err = different.MatchesOption(&prefixOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[different] query should not be supported for FieldQuery!")

	matches, err = same.Matches(original)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[same] query should not be supported for FieldQuery!")

	matches, err = same.MatchesOption(&originalOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[same] query should not be supported for FieldQuery!")

	matches, err = same.MatchesOption(&none)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[same] query should not be supported for FieldQuery!")

	matches, err = same.Matches(originalPrefix)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[same] query should not be supported for FieldQuery!")

	matches, err = same.MatchesOption(&prefixOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[same] query should not be supported for FieldQuery!")

	matches, err = empty.Matches(original)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[empty] query should not be supported for FieldQuery!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[empty] query should not be supported for FieldQuery!")

	matches, err = empty.MatchesOption(&none)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[empty] query should not be supported for FieldQuery!")

	matches, err = empty.Matches(originalPrefix)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[empty] query should not be supported for FieldQuery!")

	matches, err = empty.MatchesOption(&prefixOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[empty] query should not be supported for FieldQuery!")

	matches, err = prefix.Matches(original)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[prefix] query should not be supported for FieldQuery!")

	matches, err = prefix.MatchesOption(&originalOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[prefix] query should not be supported for FieldQuery!")

	matches, err = prefix.MatchesOption(&none)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[prefix] query should not be supported for FieldQuery!")

	matches, err = prefix.Matches(originalPrefix)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[prefix] query should not be supported for FieldQuery!")

	matches, err = prefix.MatchesOption(&prefixOption)
	assert.ErrorContains(t, err, "QueryError")
	assert.Assert(t, !matches, "Like[prefix] query should not be supported for FieldQuery!")
}

func TestLikeStringQuery(t *testing.T) {
	original := "original"
	changed := "changed"
	originalPrefix := "orig%"
	originalPattern := `o%gi_al`
	originalOption := optional.NewOption(original)
	none := optional.None[string]()
	noneCopy := optional.None[string]()
	prefixOption := optional.NewOption(originalPrefix)
	patternOption := optional.NewOption(originalPattern)

	different := smartquery.LikeString(changed)
	same := smartquery.LikeString(original)
	// Technically a user could do this so we should probably test it...
	empty := smartquery.NewStringQuery(smartquery.MatchLike, &noneCopy)
	prefix := smartquery.LikeString(originalPrefix)
	pattern := smartquery.LikeString(originalPattern)

	matches, err := different.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched different value!")

	matches, err = different.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched different value in option with Any value!")

	matches, err = different.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched option with None value!")

	matches, err = different.Matches(originalPrefix)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched prefix of original value!")

	matches, err = different.MatchesOption(&prefixOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched option with prefix of original value!")

	matches, err = different.Matches(originalPattern)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched pattern of original value!")

	matches, err = different.MatchesOption(&patternOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[different] query matched option with pattern of value!")

	matches, err = same.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[same] query did not match value when it was an exact match!")

	matches, err = same.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[same] query did not match option with Any value when it was an exact match!")

	matches, err = same.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[same] query should never match options with None value!")

	matches, err = same.Matches(originalPrefix)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[same] query matched prefix of original value!")

	matches, err = same.MatchesOption(&prefixOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[same] query matched option with prefix of original value!")

	matches, err = same.Matches(originalPattern)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[same] query matched pattern of original value!")

	matches, err = same.MatchesOption(&patternOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[same] query matched option with pattern of original value!")

	matches, err = empty.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[None] query matched any value!")

	matches, err = empty.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[None] query matches option with any value!")

	matches, err = empty.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[None] query should match options with None value!")

	matches, err = empty.Matches(originalPrefix)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[empty] query matched prefix of value!")

	matches, err = empty.MatchesOption(&prefixOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[empty] query matched option with prefix of original value!")

	matches, err = empty.Matches(originalPattern)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[empty] query matched pattern of original value!")

	matches, err = empty.MatchesOption(&patternOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[empty] query matched option with pattern of original value!")

	matches, err = prefix.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[prefix] query did not match original string!")

	matches, err = prefix.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[prefix] query did not match a string option with Any value when it should have!")

	matches, err = prefix.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[prefix] query should never match options with None value!")

	matches, err = prefix.Matches(originalPrefix)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[prefix] query did not match the same prefix!")

	matches, err = prefix.MatchesOption(&prefixOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[prefix] query matched option with the same prefix!")

	matches, err = prefix.Matches(originalPattern)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[prefix] query matched pattern of original value!")

	matches, err = prefix.MatchesOption(&patternOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[prefix] query matched option with pattern of original value!")

	matches, err = pattern.Matches(original)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[pattern] query did not match a string when it should have!")

	matches, err = pattern.MatchesOption(&originalOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[pattern] query did not match a string option with Any value when it should have!")

	matches, err = pattern.MatchesOption(&none)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[pattern] query should never match options with None value!")

	matches, err = pattern.Matches(originalPrefix)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[pattern] query pattern matched the prefix of the original value!")

	matches, err = pattern.MatchesOption(&prefixOption)
	assert.NilError(t, err)
	assert.Assert(t, !matches, "Like[pattern] query pattern matched the option with the prefix of original value!")

	matches, err = pattern.Matches(originalPattern)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[pattern] query did not match the same pattern!")

	matches, err = pattern.MatchesOption(&patternOption)
	assert.NilError(t, err)
	assert.Assert(t, matches, "Like[pattern] query did not match option with the same pattern!")
}
