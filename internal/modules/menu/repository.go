package menu

import (
	"context"
	"errors"
	"fmt"
	"kswi-backend/internal/model"
	"sort"

	"gorm.io/gorm"
)

var (
	ErrMenuNotFound      = errors.New("menu not found")
	ErrMenuHasChildren   = errors.New("cannot delete menu that has children")
	ErrCircularReference = errors.New("circular reference detected in menu hierarchy")
)

type Repository interface {
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	GetMenuByID(ctx context.Context, id uint) (*MenuResponse, error)
	CreateMenu(ctx context.Context, input CreateMenuInput) (*MenuResponse, error)
	UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) (*MenuResponse, error)
	DeleteMenu(ctx context.Context, id uint) error
	GetMenusByParentID(ctx context.Context, parentID uint) ([]MenuResponse, error)
	ValidateMenuHierarchy(ctx context.Context, menuID, newParentID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetMenuTree returns the complete active menu tree
func (r *repository) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	var menuData []model.Menu

	err := r.db.WithContext(ctx).
		Select("id, parent_id, sort, name, route, icon, is_active, created_at, updated_at").
		Where("is_active = ?", true).
		Order("sort ASC").
		Find(&menuData).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu data: %w", err)
	}

	if len(menuData) == 0 {
		return []MenuResponse{}, nil
	}

	return buildMenuTree(menuData), nil
}

// GetMenuByID returns a single menu by ID
func (r *repository) GetMenuByID(ctx context.Context, id uint) (*MenuResponse, error) {
	var menu model.Menu

	err := r.db.WithContext(ctx).
		Select("id, parent_id, sort, name, route, icon, is_active, created_at, updated_at").
		Where("id = ? AND is_active = ?", id, true).
		First(&menu).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to fetch menu: %w", err)
	}

	response := MenuResponse{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Sort:      menu.Sort,
		Name:      menu.Name,
		Route:     menu.Route,
		Icon:      menu.Icon,
		IsActive:  menu.IsActive,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
		Children:  []MenuResponse{},
	}

	return &response, nil
}

// GetMenusByParentID returns all menus with a specific parent ID
func (r *repository) GetMenusByParentID(ctx context.Context, parentID uint) ([]MenuResponse, error) {
	var menuData []model.Menu

	err := r.db.WithContext(ctx).
		Select("id, parent_id, sort, name, route, icon, is_active, created_at, updated_at").
		Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort ASC").
		Find(&menuData).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch menus by parent ID: %w", err)
	}

	var result []MenuResponse
	for _, menu := range menuData {
		result = append(result, MenuResponse{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Sort:      menu.Sort,
			Name:      menu.Name,
			Route:     menu.Route,
			Icon:      menu.Icon,
			IsActive:  menu.IsActive,
			CreatedAt: menu.CreatedAt,
			UpdatedAt: menu.UpdatedAt,
			Children:  []MenuResponse{},
		})
	}

	return result, nil
}

// CreateMenu creates a new menu with validation
func (r *repository) CreateMenu(ctx context.Context, input CreateMenuInput) (*MenuResponse, error) {
	// Validate parent exists if specified
	parentID := derefOrZero(input.ParentID, 0)
	if parentID > 0 {
		var parentExists bool
		err := r.db.WithContext(ctx).
			Model(&model.Menu{}).
			Select("1").
			Where("id = ? AND is_active = ?", parentID, true).
			Limit(1).
			Find(&parentExists).Error

		if err != nil {
			return nil, fmt.Errorf("failed to validate parent menu: %w", err)
		}
		if !parentExists {
			return nil, fmt.Errorf("parent menu with ID %d not found", parentID)
		}
	}

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	menu := model.Menu{
		ParentID:     parentID,
		Code:         input.Code,
		ParentCode:   nil,
		Sort:         input.Sort,
		Name:         input.Name,
		Route:        input.Route,
		Icon:         input.Icon,
		IsActive:     input.IsActive,
		PermissionID: nil,
	}

	if err := tx.Create(&menu).Error; err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created menu
	response := MenuResponse{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Sort:      menu.Sort,
		Name:      menu.Name,
		Route:     menu.Route,
		Icon:      menu.Icon,
		IsActive:  menu.IsActive,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
		Children:  []MenuResponse{},
	}

	return &response, nil
}

// ValidateMenuHierarchy checks for circular references
func (r *repository) ValidateMenuHierarchy(ctx context.Context, menuID, newParentID uint) error {
	if menuID == newParentID {
		return ErrCircularReference
	}

	if newParentID == 0 {
		return nil // Root level is always valid
	}

	// Check if newParentID is a descendant of menuID
	currentID := newParentID
	visited := make(map[uint]bool)

	for currentID != 0 {
		if visited[currentID] {
			return ErrCircularReference
		}
		visited[currentID] = true

		if currentID == menuID {
			return ErrCircularReference
		}

		var parentID uint
		err := r.db.WithContext(ctx).
			Model(&model.Menu{}).
			Select("parent_id").
			Where("id = ? AND is_active = ?", currentID, true).
			Scan(&parentID).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return fmt.Errorf("failed to validate hierarchy: %w", err)
		}

		currentID = parentID
	}

	return nil
}

// UpdateMenu updates an existing menu with validation
func (r *repository) UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) (*MenuResponse, error) {
	// Check if menu exists
	var existingMenu model.Menu
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&existingMenu).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}

	// Validate hierarchy if parent is being changed
	if input.ParentID != nil {
		newParentID := *input.ParentID
		if newParentID != existingMenu.ParentID {
			if err := r.ValidateMenuHierarchy(ctx, id, newParentID); err != nil {
				return nil, err
			}

			// Validate new parent exists if not root
			if newParentID > 0 {
				var parentExists bool
				err := r.db.WithContext(ctx).
					Model(&model.Menu{}).
					Select("1").
					Where("id = ? AND is_active = ?", newParentID, true).
					Limit(1).
					Find(&parentExists).Error

				if err != nil {
					return nil, fmt.Errorf("failed to validate parent menu: %w", err)
				}
				if !parentExists {
					return nil, fmt.Errorf("parent menu with ID %d not found", newParentID)
				}
			}
		}
	}

	// Build updates map
	updates := buildUpdateMap(input)
	if len(updates) == 0 {
		// Return current menu if no updates
		response := MenuResponse{
			ID:        existingMenu.ID,
			ParentID:  existingMenu.ParentID,
			Sort:      existingMenu.Sort,
			Name:      existingMenu.Name,
			Route:     existingMenu.Route,
			Icon:      existingMenu.Icon,
			IsActive:  existingMenu.IsActive,
			CreatedAt: existingMenu.CreatedAt,
			UpdatedAt: existingMenu.UpdatedAt,
			Children:  []MenuResponse{},
		}
		return &response, nil
	}

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Perform update
	result := tx.Model(&model.Menu{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update menu: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, ErrMenuNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch updated menu
	var updatedMenu model.Menu
	err = r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&updatedMenu).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated menu: %w", err)
	}

	response := MenuResponse{
		ID:        updatedMenu.ID,
		ParentID:  updatedMenu.ParentID,
		Sort:      updatedMenu.Sort,
		Name:      updatedMenu.Name,
		Route:     updatedMenu.Route,
		Icon:      updatedMenu.Icon,
		IsActive:  updatedMenu.IsActive,
		CreatedAt: updatedMenu.CreatedAt,
		UpdatedAt: updatedMenu.UpdatedAt,
		Children:  []MenuResponse{},
	}

	return &response, nil
}

// DeleteMenu performs soft delete with validation
func (r *repository) DeleteMenu(ctx context.Context, id uint) error {
	// Check if menu has children
	var childCount int64
	err := r.db.WithContext(ctx).
		Model(&model.Menu{}).
		Where("parent_id = ? AND deleted_at IS NULL", id).
		Count(&childCount).Error

	if err != nil {
		return fmt.Errorf("failed to check for children: %w", err)
	}

	if childCount > 0 {
		return ErrMenuHasChildren
	}

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Perform soft delete
	result := tx.Where("id = ?", id).Delete(&model.Menu{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete menu: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrMenuNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// buildMenuTree efficiently builds a hierarchical menu tree
func buildMenuTree(allMenus []model.Menu) []MenuResponse {
	if len(allMenus) == 0 {
		return []MenuResponse{}
	}

	// Create a map for O(1) lookups
	nodeMap := make(map[uint]*MenuResponse, len(allMenus))
	var rootNodes []MenuResponse

	// First pass: create all nodes
	for _, menu := range allMenus {
		node := &MenuResponse{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Sort:      menu.Sort,
			Name:      menu.Name,
			Route:     menu.Route,
			Icon:      menu.Icon,
			IsActive:  menu.IsActive,
			CreatedAt: menu.CreatedAt,
			UpdatedAt: menu.UpdatedAt,
			Children:  []MenuResponse{}, // Initialize empty slice
		}
		nodeMap[menu.ID] = node
	}

	// Second pass: build relationships and collect root nodes
	for _, menu := range allMenus {
		node := nodeMap[menu.ID]

		if menu.ParentID == 0 {
			// Root node
			rootNodes = append(rootNodes, *node)
		} else {
			// Child node - add to parent's children
			if parent, exists := nodeMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	// Sort root nodes
	sort.Slice(rootNodes, func(i, j int) bool {
		return rootNodes[i].Sort < rootNodes[j].Sort
	})

	// Sort children recursively
	sortChildrenRecursively(rootNodes)

	return rootNodes
}

// sortChildrenRecursively sorts children at all levels
func sortChildrenRecursively(nodes []MenuResponse) {
	for i := range nodes {
		if len(nodes[i].Children) > 0 {
			sort.Slice(nodes[i].Children, func(j, k int) bool {
				return nodes[i].Children[j].Sort < nodes[i].Children[k].Sort
			})
			sortChildrenRecursively(nodes[i].Children)
		}
	}
}

// buildUpdateMap creates the updates map for GORM
func buildUpdateMap(input UpdateMenuInput) map[string]interface{} {
	updates := make(map[string]interface{})

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

	return updates
}

// derefOrZero returns the value of pointer or default if nil
func derefOrZero[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
