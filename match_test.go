package query_test

import (
	"testing"

	"github.com/brnsampson/smartquery"
	"gotest.tools/v3/assert"
)

func TestMatcher(t *testing.T) {
	original := 47
	same := smartquery.Exact(47)
	different := smartquery.Exact(42)

	m1 := smartquery.NewValueMatch[int](original, &same)
	m2 := smartquery.NewValueMatch[int](original, &different)

	matched, err := m1.Match()
	assert.NilError(t, err)
	assert.Assert(t, matched, "An Exact query for 47 did not match 47!")

	matched, err = m2.Match()
	assert.NilError(t, err)
	assert.Assert(t, !matched, "An Exact query for 42 matched 47!")

	all := []smartquery.Matcher{m1, m2}
	matched, err = smartquery.MatchAll(all)
	assert.NilError(t, err)
	assert.Assert(t, !matched, "MatchAll with one match that should have failed returned true!")
}
