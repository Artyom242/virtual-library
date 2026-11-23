package pagination

const (
	DefaultLimit = 10
	MaxLimit     = 100
	DefaultPage  = 1
)

type Pagination struct {
	limit  int
	offset int
}

func NewPagination(page, limit *int) *Pagination {
	currentPage := DefaultPage
	currentLimit := DefaultLimit

	if page != nil {
		currentPage = *page
	}

	if limit != nil {
		currentLimit = *limit
	}

	if currentPage < 1 {
		currentPage = DefaultPage
	}

	if currentLimit < 1 {
		currentLimit = DefaultLimit
	}

	if currentLimit > MaxLimit {
		currentLimit = MaxLimit
	}

	return &Pagination{
		limit:  currentLimit,
		offset: newOffset(currentPage, currentLimit),
	}
}

func newOffset(page, limit int) int {
	return (page - 1) * limit
}

func (p Pagination) Limit() int {
	return p.limit
}
func (p Pagination) Offset() int {
	return p.offset
}
