package models

import "strconv"

type Filter struct {
	Name, Value string
}

type Pagination struct {
	limit   int
	offset  int
	filters []Filter
}

func NewPagination(filters ...Filter) *Pagination {
	return &Pagination{
		filters: filters,
	}
}

func (p *Pagination) Limit() int {
	return p.limit
}

func (p *Pagination) Offset() int {
	return p.offset
}

func (p *Pagination) Filter(name string) (Filter, bool) {
	for _, f := range p.filters {
		if f.Name == name {
			return f, true
		}
	}
	return Filter{}, false
}

func (p *Pagination) ToPairs() []string {
	var pairs []string

	if p.limit > 0 {
		pairs = append(pairs, "limit", strconv.Itoa(p.limit))
	}

	if p.offset > 0 {
		pairs = append(pairs, "offset", strconv.Itoa(p.offset))
	}

	for _, f := range p.filters {
		if f.Value != "" {
			pairs = append(pairs, f.Name, f.Value)
		}
	}

	return pairs
}
