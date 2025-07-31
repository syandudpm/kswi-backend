package menu

import (
	"time"
)

// MenuResponse represents a menu item in API responses
type MenuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Sort      int            `json:"sort"`
	Name      string         `json:"name"`
	Route     *string        `json:"route,omitempty"`
	Icon      *string        `json:"icon,omitempty"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Children  []MenuResponse `json:"children"`

	// Additional computed fields (added by service layer)
	HasAccess  *bool `json:"has_access,omitempty"`
	Depth      *int  `json:"depth,omitempty"`
	ChildCount *int  `json:"child_count,omitempty"`
}

// CreateMenuInput represents input for creating a new menu
type CreateMenuInput struct {
	ParentID *uint   `json:"parent_id,omitempty" validate:"omitempty,min=1"`
	Code     *string `json:"code,omitempty" validate:"omitempty,max=200"`
	Sort     int     `json:"sort" validate:"min=0,max=9999"`
	Name     string  `json:"name" validate:"required,max=200"`
	Route    *string `json:"route,omitempty" validate:"omitempty,max=200"`
	Icon     *string `json:"icon,omitempty" validate:"omitempty,max=1000"`
	IsActive bool    `json:"is_active"`
}

// UpdateMenuInput represents input for updating an existing menu
type UpdateMenuInput struct {
	ParentID *uint   `json:"parent_id,omitempty" validate:"omitempty,min=0"`
	Code     *string `json:"code,omitempty" validate:"omitempty,max=200"`
	Sort     *int    `json:"sort,omitempty" validate:"omitempty,min=0,max=9999"`
	Name     *string `json:"name,omitempty" validate:"omitempty,required,max=200"`
	Route    *string `json:"route,omitempty" validate:"omitempty,max=200"`
	Icon     *string `json:"icon,omitempty" validate:"omitempty,max=1000"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// MenuOrderInput represents input for reordering menus
type MenuOrderInput struct {
	ID   uint `json:"id" validate:"required,min=1"`
	Sort int  `json:"sort" validate:"min=0,max=9999"`
}

// MenuSearchInput represents input for searching menus
type MenuSearchInput struct {
	Query    string `json:"query" validate:"required,min=1,max=100"`
	Page     int    `json:"page,omitempty" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size,omitempty" validate:"omitempty,min=1,max=100"`
}

// MenuBulkStatusInput represents input for bulk status updates
type MenuBulkStatusInput struct {
	MenuIDs  []uint `json:"menu_ids" validate:"required,min=1,dive,min=1"`
	IsActive bool   `json:"is_active"`
}

// MenuStatistics represents menu usage statistics
type MenuStatistics struct {
	TotalMenus    int       `json:"total_menus"`
	ActiveMenus   int       `json:"active_menus"`
	InactiveMenus int       `json:"inactive_menus"`
	RootMenus     int       `json:"root_menus"`
	MaxDepth      int       `json:"max_depth"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// MenuTreeOptions represents options for retrieving menu tree
type MenuTreeOptions struct {
	IncludeInactive bool  `json:"include_inactive,omitempty"`
	MaxDepth        *int  `json:"max_depth,omitempty"`
	UserID          *uint `json:"user_id,omitempty"`       // For permission filtering
	StartFromID     *uint `json:"start_from_id,omitempty"` // Start tree from specific menu
}

// MenuPathResponse represents a breadcrumb path
type MenuPathResponse struct {
	Path  []MenuResponse `json:"path"`
	Total int            `json:"total"`
}

// MenuSearchResponse represents search results
type MenuSearchResponse struct {
	Results []MenuResponse `json:"results"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	Pages   int            `json:"pages"`
}

// MenuValidationError represents validation errors
type MenuValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// MenuErrorResponse represents error response
type MenuErrorResponse struct {
	Error       string                 `json:"error"`
	Message     string                 `json:"message"`
	Code        string                 `json:"code,omitempty"`
	Validations []MenuValidationError  `json:"validations,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// HTTP Request/Response structures for handlers

// CreateMenuRequest represents HTTP request for creating menu
type CreateMenuRequest struct {
	CreateMenuInput
}

// UpdateMenuRequest represents HTTP request for updating menu
type UpdateMenuRequest struct {
	UpdateMenuInput
}

// ReorderMenusRequest represents HTTP request for reordering menus
type ReorderMenusRequest struct {
	Orders []MenuOrderInput `json:"orders" validate:"required,min=1,dive"`
}

// BulkStatusRequest represents HTTP request for bulk status update
type BulkStatusRequest struct {
	MenuBulkStatusInput
}

// GetMenuTreeRequest represents HTTP request for getting menu tree
type GetMenuTreeRequest struct {
	IncludeInactive *bool `query:"include_inactive"`
	MaxDepth        *int  `query:"max_depth" validate:"omitempty,min=1,max=10"`
	UserID          *uint `query:"user_id" validate:"omitempty,min=1"`
	StartFromID     *uint `query:"start_from_id" validate:"omitempty,min=1"`
}

// SearchMenusRequest represents HTTP request for searching menus
type SearchMenusRequest struct {
	Query    string `query:"q" validate:"required,min=1,max=100"`
	Page     *int   `query:"page" validate:"omitempty,min=1"`
	PageSize *int   `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// Standard API Response structures

// SuccessResponse represents successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ListResponse represents paginated list response
type ListResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// Meta represents pagination metadata
type Meta struct {
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	Pages    int  `json:"pages"`
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
}

// Response helper functions

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message ...string) *SuccessResponse {
	resp := &SuccessResponse{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

// NewListResponse creates a paginated list response
func NewListResponse(data interface{}, page, pageSize, total int) *ListResponse {
	pages := (total + pageSize - 1) / pageSize
	if pages < 1 {
		pages = 1
	}

	return &ListResponse{
		Success: true,
		Data:    data,
		Meta: Meta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    pages,
			HasNext:  page < pages,
			HasPrev:  page > 1,
		},
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err error, code ...string) *MenuErrorResponse {
	resp := &MenuErrorResponse{
		Error:   "Error",
		Message: err.Error(),
	}
	if len(code) > 0 {
		resp.Code = code[0]
	}
	return resp
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(validations []MenuValidationError) *MenuErrorResponse {
	return &MenuErrorResponse{
		Error:       "Validation Error",
		Message:     "Input validation failed",
		Code:        "VALIDATION_ERROR",
		Validations: validations,
	}
}

// Utility functions for DTOs

// SetDefaults sets default values for CreateMenuInput
func (c *CreateMenuInput) SetDefaults() {
	// Set default values if not provided
	if c.Code == nil || *c.Code == "" {
		// Will be generated by service
	}
}

// SetDefaults sets default values for MenuSearchInput
func (s *MenuSearchInput) SetDefaults() {
	if s.Page == 0 {
		s.Page = 1
	}
	if s.PageSize == 0 {
		s.PageSize = 20
	}
}

// Validate checks if CreateMenuInput is valid
func (c *CreateMenuInput) Validate() []MenuValidationError {
	var errors []MenuValidationError

	if c.Name == "" {
		errors = append(errors, MenuValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if len(c.Name) > 200 {
		errors = append(errors, MenuValidationError{
			Field:   "name",
			Message: "Name must not exceed 200 characters",
			Value:   len(c.Name),
		})
	}

	if c.Sort < 0 || c.Sort > 9999 {
		errors = append(errors, MenuValidationError{
			Field:   "sort",
			Message: "Sort must be between 0 and 9999",
			Value:   c.Sort,
		})
	}

	if c.Route != nil && len(*c.Route) > 200 {
		errors = append(errors, MenuValidationError{
			Field:   "route",
			Message: "Route must not exceed 200 characters",
			Value:   len(*c.Route),
		})
	}

	return errors
}

// HasChanges checks if UpdateMenuInput has any changes
func (u *UpdateMenuInput) HasChanges() bool {
	return u.ParentID != nil ||
		u.Code != nil ||
		u.Sort != nil ||
		u.Name != nil ||
		u.Route != nil ||
		u.Icon != nil ||
		u.IsActive != nil
}

// AddChild adds a child menu to the response (helper function)
func (m *MenuResponse) AddChild(child MenuResponse) {
	if m.Children == nil {
		m.Children = []MenuResponse{}
	}
	m.Children = append(m.Children, child)
}

// GetChildCount returns the number of direct children
func (m *MenuResponse) GetChildCount() int {
	return len(m.Children)
}

// GetTotalDescendants returns the total number of descendants (recursive)
func (m *MenuResponse) GetTotalDescendants() int {
	total := len(m.Children)
	for _, child := range m.Children {
		total += child.GetTotalDescendants()
	}
	return total
}

// IsRoot checks if the menu is a root menu
func (m *MenuResponse) IsRoot() bool {
	return m.ParentID == 0
}

// HasChildren checks if the menu has children
func (m *MenuResponse) HasChildren() bool {
	return len(m.Children) > 0
}
