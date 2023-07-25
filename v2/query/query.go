package query

import (
	"github.com/tauraamui/kvs/v2"
	"github.com/tauraamui/kvs/v2/storage"
)

type Query struct {
	filters []Filter
}

type operator int64

const (
	undefined operator = iota
	equal
	lessthan
)

func (op operator) String() string {
	switch op {
	case equal:
		return "equal"
	default:
		return "undefined"
	}
}

type Filter struct {
	q         *Query
	fieldName string
	op        operator
	value     any
}

func New() *Query {
	return &Query{}
}

func Run[T storage.Value](s storage.Store, q *Query) ([]T, error) {
	return storage.LoadAllWithOperators[T](s, kvs.RootOwner{}, func(e kvs.Entry) bool {
		for _, filter := range q.filters {
			if filter.op == equal {
				return kvs.CompareBytesToAny(e.Data, filter.value)
			}
		}

		return len(q.filters) > 0
	})
}

func (q *Query) Filter(fieldName string) *Filter {
	q = q.clone()
	filter := Filter{q: q}
	q.filters = append(q.filters, filter)
	return &q.filters[len(q.filters)-1]
}

func (f *Filter) Eq(value any) *Query {
	f.value = value
	f.op = equal
	return f.q
}

func (f *Filter) Lt(value any) *Query {
	f.value = value
	f.op = lessthan
	return f.q
}

func (q *Query) clone() *Query {
	x := *q
	// Copy the contents of the slice-typed fields to a new backing store.
	if len(q.filters) > 0 {
		x.filters = make([]Filter, len(q.filters))
		copy(x.filters, q.filters)
	}
	return &x
}
