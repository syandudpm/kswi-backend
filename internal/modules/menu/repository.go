package menu

import (
	"context"
	"fmt"
	"sort"

	"kswi-backend/internal/config"
	"kswi-backend/internal/models"
	"kswi-backend/internal/shared/pagination"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, params pagination.Params) ([]models.User, *pagination.Meta, error)
	Count(ctx context.Context) (int64, error)

	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository() Repository {
	return &repository{
		db: config.GetDB(),
	}
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *repository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with pagination
func (r *repository) List(ctx context.Context, params pagination.Params) ([]models.User, *pagination.Meta, error) {
	var users []models.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate pagination
	meta := pagination.Calculate(params, total)

	// Query with pagination
	query := r.db.WithContext(ctx).
		Offset(meta.Offset).
		Limit(meta.Limit)

	// Add sorting
	if params.Sort != "" {
		query = query.Order(params.Sort + " " + params.Order)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, meta, nil
}

// Count returns the total number of users
func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

type menuData struct {
	ID       uint    `gorm:"column:id"`
	ParentID uint    `gorm:"column:parent_id"` // Nullable for root items
	Sort     int     `gorm:"column:sort"`
	Name     string  `gorm:"column:name"`
	Route    *string `gorm:"column:route"`
	Icon     *string `gorm:"column:icon"`
}

func (r *repository) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	// Step 1: Fetch ALL active menus in a single query
	var allMenus []menuData
	err := r.db.WithContext(ctx).
		Table("golang.menus").
		Select("id, parent_id, sort, name, route, icon").
		Where("is_active = 1").
		Order("sort ASC"). // Sort all items by sort
		Find(&allMenus).Error

	if err != nil {
		return nil, err
	}

	// Step 2: Build tree structure recursively starting from root (parentID = 0)
	return r.buildMenuTree(allMenus, 0), nil
}

func (r *repository) buildMenuTree(allMenus []menuData, parentID uint) []MenuResponse {
	var menus []MenuResponse

	// Find all menus with the given parentID
	for _, menu := range allMenus {
		if menu.ParentID == parentID {

			// RECURSIVE CALL: Build children for this menu
			children := r.buildMenuTree(allMenus, menu.ID) // <-- HERE'S THE RECURSION!

			// Ensure children is empty array instead of nil
			// if children == nil {
			// 	children = []MenuResponse{}
			// }

			menuResponse := MenuResponse{
				ID:       menu.ID,
				Sort:     menu.Sort,
				Name:     menu.Name,
				Route:    menu.Route,
				Icon:     menu.Icon,
				Children: children, // Assign the recursively built children
			}
			menus = append(menus, menuResponse)
		}
	}

	// Sort menus by sort
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Sort < menus[j].Sort
	})

	return menus
}
