package menu

import "time"

// MenuResponse is the API response DTO
type MenuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Sort      int            `json:"sort"`
	Name      string         `json:"name"`
	Route     *string        `json:"route"`
	Icon      *string        `json:"icon"`
	IsActive  bool           `json:"is_active"`
	Children  []MenuResponse `json:"children"`
	CreatedAt time.Time      `json:"created_at"`
}

// CreateMenuInput is for creating a new menu
type CreateMenuInput struct {
	ParentID *uint   `json:"parent_id" binding:"omitempty"`
	Code     *string `json:"code" binding:"omitempty,max=200"`
	Sort     int     `json:"sort" binding:"gte=0"`
	Name     string  `json:"name" binding:"required,max=200"`
	Route    *string `json:"route" binding:"omitempty,max=200"`
	Icon     *string `json:"icon" binding:"omitempty,max=1000"`
	IsActive bool    `json:"is_active"`
}

// UpdateMenuInput – only editable fields
type UpdateMenuInput struct {
	ParentID *uint   `json:"parent_id,omitempty"`
	Code     *string `json:"code,omitempty" binding:"omitempty,max=200"`
	Sort     *int    `json:"sort,omitempty" binding:"omitempty,gte=0"`
	Name     *string `json:"name,omitempty" binding:"omitempty,max=200"`
	Route    *string `json:"route,omitempty" binding:"omitempty,max=200"`
	Icon     *string `json:"icon,omitempty" binding:"omitempty,max=1000"`
	IsActive *bool   `json:"is_active,omitempty"`
}
