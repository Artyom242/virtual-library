package httpapi

import "virtual-library/catalog-service/internal/domain"

func parsePagination(page, limit *int) (int, int, error) {
	p := 1
	if page != nil {
		p = *page
	}

	l := 10
	if limit != nil {
		l = *limit
	}

	if p < 1 || l < 1 || l > 100 {
		return 0, 0, domain.ErrInvalidPagination
	}

	return p, l, nil
}
