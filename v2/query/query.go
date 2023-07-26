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

func (f Filter) cmp(d []byte) bool {
	return kvs.CompareBytesToAny(d, f.value)
}

func New() *Query {
	return &Query{}
}

func Run[T storage.Value](s storage.Store, q *Query) ([]T, error) {
	return storage.LoadAllWithEvaluator[T](s, kvs.RootOwner{}, func(e kvs.Entry) bool {
		if q == nil {
			return true
		}

		matching := false
		for _, filter := range q.filters {
			if filter.fieldName == e.ColumnName {
				if filter.op == equal {
					matching = filter.cmp(e.Data)
					if !matching {
						return false
					}
				}
			}
		}

		return matching || len(q.filters) == 0
	})
}

func (q *Query) Filter(fieldName string) *Filter {
	q = q.clone()
	filter := Filter{q: q, fieldName: fieldName}
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
