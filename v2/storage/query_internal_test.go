package storage

import (
	"testing"

	"github.com/matryer/is"
)

func TestOperatorString(t *testing.T) {
	is := is.New(t)
	eq := equal
	un := undefined

	is.Equal(eq.String(), "equal")
	is.Equal(un.String(), "undefined")

}

func TestQueryFilterSetsOperatorAndValue(t *testing.T) {
	is := is.New(t)
	q := NewQuery().Filter("color").Eq("yellow")

	is.True(len(q.filters) == 1)
	is.Equal(q.filters[0].op, equal)
	is.Equal(q.filters[0].value, "yellow")
}

func TestQueryFilterSubsequentOperatorAndValueOverwritesPrevious(t *testing.T) {
	is := is.New(t)
	q := NewQuery()
	filter := q.Filter("color")
	q = filter.Eq("yellow")
	q = filter.Eq("blue")

	is.True(len(q.filters) == 1)
	is.Equal(q.filters[0].op, equal)
	is.Equal(q.filters[0].value, "blue")
}
