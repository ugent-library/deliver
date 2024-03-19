package models

import (
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const defaultOffset = 0
const defaultLimit = 20
const defaultSort = "default"

type Filter struct {
	Name, Value string
}

type Pagination struct {
	offset  int
	limit   int
	sort    string
	total   int
	filters []Filter
}

func NewPagination(offset int, limit int, sort string, filters ...Filter) *Pagination {
	if offset < 0 {
		offset = defaultOffset
	}

	if limit <= 0 {
		limit = defaultLimit
	}

	return &Pagination{
		offset:  offset,
		limit:   limit,
		total:   -1,
		sort:    sort,
		filters: filters,
	}
}

func (p *Pagination) Offset() int {
	return p.offset
}

func (p *Pagination) Limit() int {
	return p.limit
}

func (p *Pagination) Total() int {
	return p.total
}

func (p *Pagination) SetTotal(total int) {
	p.total = total
}

func (p *Pagination) StartOfPage() int {
	return p.Offset() + 1
}

func (p *Pagination) EndOfPage() int {
	if p.Total() >= 0 {
		return min(p.Total(), p.Offset()+p.Limit())
	}

	return p.Offset() + p.Limit()
}

func (p *Pagination) CurrentPage() int {
	if p.Total() < 0 {
		return -1
	}

	return int(math.Floor(float64(p.Offset())/float64(p.Limit()))) + 1
}

func (p *Pagination) NumberOfPages() int {
	if p.Total() < 0 {
		return -1
	}

	return int(math.Ceil(float64(p.Total()) / float64(p.Limit())))
}

func (p *Pagination) PageOffset(page int) int {
	return (page - 1) * p.Limit()
}

func (p *Pagination) VisiblePages() []int {
	if p.total < 0 {
		return []int{1}
	}

	// Initialize map with each possible page number
	pages := make(map[int]bool, p.NumberOfPages())

	// Activate pages 1, current - 1, current, current + 1 and end
	pages[1] = true
	pages[p.CurrentPage()-1] = true
	pages[p.CurrentPage()] = true
	pages[p.CurrentPage()+1] = true
	pages[p.NumberOfPages()] = true

	// Delete pages outside boundaries
	for page := range pages {
		if page != 1 && (page < 1 || page > p.NumberOfPages()) {
			delete(pages, page)
		}
	}

	// Activate single-gap-pages ("[4] […] [6]" => "[4] [5] [6]")
	for i := 1; i < p.NumberOfPages(); i++ {
		if !pages[i] && pages[i-1] && pages[i+1] {
			pages[i] = true
		}
	}

	// Create sorted list of visible pages
	visiblePages := make([]int, 0, p.NumberOfPages())
	for page, isVisible := range pages {
		if isVisible {
			visiblePages = append(visiblePages, page)
		}
	}
	sort.Ints(visiblePages)

	// Add ellipsis .pager-items in between gaps
	for i := 0; i < len(visiblePages); i++ {
		if i >= 1 && visiblePages[i] != visiblePages[i-1]+1 {
			visiblePages = slices.Insert(visiblePages, i, -1)
			i++
		}
	}

	return visiblePages
}

func (p *Pagination) Sort() string {
	return p.sort
}

func (p *Pagination) Filter(name string) (Filter, bool) {
	for _, f := range p.filters {
		if f.Name == name {
			return f, true
		}
	}
	return Filter{}, false
}

func (p *Pagination) FilterValue(name string) string {
	filter, ok := p.Filter(name)
	if !ok {
		return ""
	}

	return filter.Value
}

func (p *Pagination) ToPairs() []string {
	var pairs []string

	for _, f := range p.filters {
		if f.Value != "" {
			pairs = append(pairs, f.Name, f.Value)
		}
	}

	if p.sort != "" && p.sort != defaultSort {
		pairs = append(pairs, "sort", p.sort)
	}

	if p.offset != defaultOffset {
		pairs = append(pairs, "offset", strconv.Itoa(p.offset))
	}

	if p.limit != defaultLimit {
		pairs = append(pairs, "limit", strconv.Itoa(p.limit))
	}

	return pairs
}

func (p *Pagination) PaginationString() string {
	var sb strings.Builder

	if p.total != 0 {
		sb.WriteString(strconv.Itoa(p.StartOfPage()))
		sb.WriteString("-")
		sb.WriteString(strconv.Itoa(p.EndOfPage()))
		sb.WriteString(" of ")
	}

	if p.total >= 0 {
		sb.WriteString(strconv.Itoa(p.total))
	} else {
		sb.WriteString("???")
	}

	return sb.String()
}
