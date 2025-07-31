package pagination

// PaginationRequest represents common pagination request parameters
type PaginationRequest struct {
	Page     int      `json:"page" binding:"required,min=1"`
	PerPage  int      `json:"per_page" binding:"required,min=1,max=100"`
	Search   string   `json:"search"`
	SortBy   string   `json:"sort_by"`
	SortDesc bool     `json:"sort_desc"`
	Filters  *Filters `json:"filters"`
}

// QueryParams returns the query parameters for SQL
func (r *PaginationRequest) QueryParams() (offset int, limit int) {
	offset = (r.Page - 1) * r.PerPage
	limit = r.PerPage
	return
}

// SortParams returns the sorting parameters
func (r *PaginationRequest) SortParams() (column string, direction string) {
	if r.SortBy == "" {
		return "", ""
	}

	direction = "ASC"
	if r.SortDesc {
		direction = "DESC"
	}

	return r.SortBy, direction
}

// WithDefaultSort sets default sort if none provided
func (r *PaginationRequest) WithDefaultSort(field string, desc bool) {
	if r.SortBy == "" {
		r.SortBy = field
		r.SortDesc = desc
	}
}

// WithMaxPerPage ensures per_page doesn't exceed maximum
func (r *PaginationRequest) WithMaxPerPage(max int) {
	if r.PerPage > max {
		r.PerPage = max
	}
}
