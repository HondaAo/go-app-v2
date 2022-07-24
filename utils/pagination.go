package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultSize = 10
)

// Pagination query params
type PaginationQuery struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

// Set page size
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// Get pagination query struct from
func GetPaginationFromCtx(c echo.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.QueryParam("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.QueryParam("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.QueryParam("orderBy"))

	return q, nil
}
