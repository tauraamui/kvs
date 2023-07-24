package storage

type Query struct {
	filters []Filter
}

type Filter struct {
	q         *Query
	fieldName string
	operator  string
	value     any
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Filter(fieldName string) *Filter {
	q = q.clone()
	filter := Filter{q: q}
	q.filters = append(q.filters, filter)
	return &q.filters[len(q.filters)-1]
}

func (f *Filter) Eq(value any) *Query {
	f.operator = "eq"
	f.value = value
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
