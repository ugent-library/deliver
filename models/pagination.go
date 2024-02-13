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
	var usedFilters []string

	if p.limit > 0 {
		usedFilters = append(usedFilters, "limit", strconv.Itoa(p.limit))
	}

	if p.offset > 0 {
		usedFilters = append(usedFilters, "offset", strconv.Itoa(p.offset))
	}

	for _, f := range p.filters {
		if f.Value != "" {
			usedFilters = append(usedFilters, f.Name, f.Value)
		}
	}

	return usedFilters
}
