package menu

import (
	"context"
	"database/sql"
	"kswi-backend/internal/model"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	CreateMenu(ctx context.Context, input CreateMenuInput) error
	UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) error
	DeleteMenu(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetMenuTree returns the active menu tree using raw SQL → DTO
func (r *repository) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	const query = `
        SELECT id, parent_id, sort, name, route, icon, is_active, created_at
        FROM menus
        WHERE is_active = true AND deleted_at IS NULL
        ORDER BY sort ASC
    `

	var rows []struct {
		ID        uint      `db:"id"`
		ParentID  uint      `db:"parent_id"`
		Sort      int       `db:"sort"`
		Name      string    `db:"name"`
		Route     *string   `db:"route"`
		Icon      *string   `db:"icon"`
		IsActive  bool      `db:"is_active"`
		CreatedAt time.Time `db:"created_at"`
	}

	err := r.db.WithContext(ctx).Raw(query).Scan(&rows).Error
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return buildMenuTree(rows), nil
}

func buildMenuTree(rows []struct {
	ID        uint
	ParentID  uint
	Sort      int
	Name      string
	Route     *string
	Icon      *string
	IsActive  bool
	CreatedAt time.Time
}) []MenuResponse {
	idMap := make(map[uint]MenuResponse)
	var result []MenuResponse

	// Initialize all nodes
	for _, row := range rows {
		idMap[row.ID] = MenuResponse{
			ID:        row.ID,
			ParentID:  row.ParentID,
			Sort:      row.Sort,
			Name:      row.Name,
			Route:     row.Route,
			Icon:      row.Icon,
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt,
			Children:  []MenuResponse{}, // never nil
		}
	}

	// Link children to parents
	for _, row := range rows {
		if row.ParentID == 0 {
			result = append(result, idMap[row.ID])
		} else {
			if parent, exists := idMap[row.ParentID]; exists {
				parent.Children = append(parent.Children, idMap[row.ID])
				idMap[row.ParentID] = parent // update parent (slice is copied)
			}
		}
	}

	// Sort by sort field
	sort.Slice(result, func(i, j int) bool {
		return result[i].Sort < result[j].Sort
	})

	return result
}

// CreateMenu uses GORM with full model
func (r *repository) CreateMenu(ctx context.Context, input CreateMenuInput) error {
	menu := model.Menu{
		ParentID:     derefOrZero(input.ParentID, 0),
		Code:         input.Code,
		ParentCode:   nil, // can be set in service if needed
		Sort:         input.Sort,
		Name:         input.Name,
		Route:        input.Route,
		Icon:         input.Icon,
		IsActive:     input.IsActive,
		PermissionID: nil,
	}

	return r.db.WithContext(ctx).Create(&menu).Error
}

// UpdateMenu – only updates allowed fields
func (r *repository) UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) error {
	updates := map[string]interface{}{}
	if input.ParentID != nil {
		updates["parent_id"] = *input.ParentID
	}
	if input.Code != nil {
		updates["code"] = input.Code
	}
	if input.Sort != nil {
		updates["sort"] = *input.Sort
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Route != nil {
		updates["route"] = input.Route
	}
	if input.Icon != nil {
		updates["icon"] = input.Icon
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) == 0 {
		return nil // nothing to update
	}

	result := r.db.WithContext(ctx).
		Model(&model.Menu{}).
		Where("id = ?", id).
		Select("parent_id", "code", "sort", "name", "route", "icon", "is_active", "updated_at").
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteMenu – soft delete
func (r *repository) DeleteMenu(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Menu{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Helper: deref pointer or return default
func derefOrZero[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
