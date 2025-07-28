package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params holds pagination parameters
type Params struct {
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order"`
	Search string `json:"search" form:"search"`
	Filter string `json:"filter" form:"filter"`
}

// Meta holds pagination metadata
type Meta struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasPrevious bool  `json:"has_previous"`
	HasNext     bool  `json:"has_next"`
	Offset      int   `json:"offset"`
}

// Response holds paginated response
type Response struct {
	Data *interface{} `json:"data"`
	Meta *Meta        `json:"meta"`
}

// DefaultLimit is the default number of items per page
const DefaultLimit = 20

// MaxLimit is the maximum number of items per page
const MaxLimit = 100

// ParseParams parses pagination parameters from gin context
func ParseParams(c *gin.Context) Params {
	params := Params{
		Page:   1,
		Limit:  DefaultLimit,
		Sort:   "",
		Order:  "DESC",
		Search: "",
		Filter: "",
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > MaxLimit {
				limit = MaxLimit
			}
			params.Limit = limit
		}
	}

	// Parse sort
	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	// Parse order
	if order := c.Query("order"); order != "" {
		if order == "ASC" || order == "DESC" {
			params.Order = order
		}
	}

	// Parse search
	params.Search = c.Query("search")

	// Parse filter
	params.Filter = c.Query("filter")

	return params
}

// Calculate calculates pagination metadata
func Calculate(params Params, total int64) *Meta {
	if params.Limit <= 0 {
		params.Limit = DefaultLimit
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Ensure page doesn't exceed total pages
	if params.Page > totalPages {
		params.Page = totalPages
	}

	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	return &Meta{
		Page:        params.Page,
		Limit:       params.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasPrevious: params.Page > 1,
		HasNext:     params.Page < totalPages,
		Offset:      offset,
	}
}

// NewResponse creates a new paginated response
func NewResponse(data interface{}, meta *Meta) *Response {
	return &Response{
		Data: &data,
		Meta: meta,
	}
}
