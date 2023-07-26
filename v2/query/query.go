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
	values    []any
}

func (f Filter) cmp(d []byte) bool {
	for _, v := range f.values {
		if kvs.CompareBytesToAny(d, v) {
			return true
		}
	}
	return false
}

func New() *Query {
	return &Query{}
}

func Run[T storage.Value](s storage.Store, q *Query) ([]T, error) {
	return storage.LoadAllWithEvaluator[T](s, kvs.RootOwner{}, func(e kvs.Entry) bool {
		if q == nil || len(q.filters) == 0 {
			return true
		}

		captured := true
		for i, filter := range q.filters {
			if i > 0 && !captured {
				return false
			}
			if filter.fieldName == e.ColumnName {
				if filter.op == equal {
					if !filter.cmp(e.Data) {
						captured = false
					}
				}
			}
		}

		return captured
	})
}

func (q *Query) Filter(fieldName string) *Filter {
	q = q.clone()
	filter := Filter{q: q, fieldName: fieldName}
	q.filters = append(q.filters, filter)
	return &q.filters[len(q.filters)-1]
}

func (f *Filter) Eq(value ...any) *Query {
	f.values = value
	f.op = equal
	return f.q
}

func (f *Filter) Lt(value ...any) *Query {
	f.values = value
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
