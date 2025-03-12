package utils

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FilterModel struct {
	Field      string `query:"field" json:"field"`
	Value      string `query:"value" json:"value"`
	Comparison string `query:"comparison" json:"comparison"`
}

const (
	defaultSize = 10
	defaultPage = 1
	maxSize     = 500
)

type Pagination struct {
	Size    int            `query:"size" json:"size,omitempty"`
	Page    int            `query:"page" json:"page,omitempty"`
	OrderBy string         `query:"orderBy" json:"orderBy,omitempty"`
	Filters []*FilterModel `query:"filters" json:"filters,omitempty"`
}

func NewPagination(size int, page int) *Pagination {
	if size > maxSize {
		size = maxSize
	}
	return &Pagination{Size: size, Page: page}
}
func NewPaginationFromQueryParams(size string, page string) *Pagination {

	p := &Pagination{Size: defaultSize, Page: defaultPage}

	if sizeNum, err := strconv.Atoi(size); err == nil && sizeNum != 0 {
		p.Page = sizeNum

		if sizeNum > maxSize {
			p.Size = maxSize
		}
	}

	if pageNum, err := strconv.Atoi(page); err == nil && pageNum != 0 {
		p.Page = pageNum
	}

	return p
}

func ToPagination(c echo.Context) (*Pagination, error) {
	q := &Pagination{}
	var page, size, orderBy string

	err := echo.QueryParamsBinder(c).
		String("size", &size).
		String("page", &page).
		String("orderBy", &orderBy).
		BindError()

	if err != nil {
		return nil, err
	}

	if err = q.SetPage(page); err != nil {
		return nil, err
	}
	if err = q.SetSize(size); err != nil {
		return nil, err
	}
	q.SetOrderBy(orderBy)

	// Handle filters separately since they need special parsing
	filters := c.QueryParams()["filters"]
	for _, filter := range filters {
		if filter == "" {
			continue
		}
		f := &FilterModel{
			Field:      c.QueryParam("field"),
			Value:      c.QueryParam("value"),
			Comparison: c.QueryParam("comparison"),
		}
		q.Filters = append(q.Filters, f)
	}

	return q, nil
}

// SetSize Set page size
func (q *Pagination) SetSize(sizeQuery string) error {
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

// SetPage Set page number
func (q *Pagination) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// SetOrderBy Set order by
func (q *Pagination) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// GetOffset Get offset
func (q *Pagination) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit Get limit
func (q *Pagination) GetLimit() int {
	return q.Size
}

// GetOrderBy Get OrderBy
func (q *Pagination) GetOrderBy() string {
	return q.OrderBy
}

// GetPage Get OrderBy
func (q *Pagination) GetPage() int {
	return q.Page
}

// GetSize Get OrderBy
func (q *Pagination) GetSize() int {
	return q.Size
}

// GetQueryString get query string
func (q *Pagination) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}
