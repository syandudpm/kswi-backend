package menu

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Service interface defines all business operations
type Service interface {
	// Read operations
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	GetMenuByID(ctx context.Context, id uint) (*MenuResponse, error)
	GetMenusByParentID(ctx context.Context, parentID uint) ([]MenuResponse, error)
	GetMenuPath(ctx context.Context, menuID uint) ([]MenuResponse, error)
	SearchMenus(ctx context.Context, query string) ([]MenuResponse, error)

	// Write operations
	CreateMenu(ctx context.Context, input CreateMenuInput) (*MenuResponse, error)
	UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) (*MenuResponse, error)
	DeleteMenu(ctx context.Context, id uint) error

	// Bulk operations
	ReorderMenus(ctx context.Context, orders []MenuOrderInput) error
	BulkUpdateMenuStatus(ctx context.Context, ids []uint, isActive bool) error

	// Utility operations
	ValidateMenuAccess(ctx context.Context, userID uint, menuID uint) (bool, error)
	GetMenuStatistics(ctx context.Context) (*MenuStatistics, error)
}

// Service errors
var (
	ErrInvalidInput      = errors.New("invalid input provided")
	ErrDuplicateMenuName = errors.New("menu name already exists at this level")
	ErrMaxDepthExceeded  = errors.New("maximum menu depth exceeded")
	ErrInvalidSortOrder  = errors.New("invalid sort order")
)

// Service constants
const (
	MaxMenuDepth    = 5
	MaxMenuNameLen  = 200
	MaxRouteLen     = 200
	MinSortValue    = 0
	MaxSortValue    = 9999
	DefaultPageSize = 20
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetMenuTree returns the complete menu hierarchy
func (s *service) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	menus, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu tree: %w", err)
	}

	// Add additional business logic if needed
	s.enrichMenuData(menus)

	return menus, nil
}

// GetMenuByID returns a single menu with validation
func (s *service) GetMenuByID(ctx context.Context, id uint) (*MenuResponse, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	menu, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMenuNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("failed to get menu by ID: %w", err)
	}

	return menu, nil
}

// GetMenusByParentID returns all menus under a specific parent
func (s *service) GetMenusByParentID(ctx context.Context, parentID uint) ([]MenuResponse, error) {
	menus, err := s.repo.GetMenusByParentID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menus by parent ID: %w", err)
	}

	return menus, nil
}

// GetMenuPath returns the breadcrumb path to a menu
func (s *service) GetMenuPath(ctx context.Context, menuID uint) ([]MenuResponse, error) {
	if menuID == 0 {
		return []MenuResponse{}, nil
	}

	var path []MenuResponse
	currentID := menuID

	// Build path by traversing up the hierarchy
	for currentID != 0 {
		menu, err := s.repo.GetMenuByID(ctx, currentID)
		if err != nil {
			if errors.Is(err, ErrMenuNotFound) {
				break
			}
			return nil, fmt.Errorf("failed to build menu path: %w", err)
		}

		// Prepend to path (so root comes first)
		path = append([]MenuResponse{*menu}, path...)
		currentID = menu.ParentID
	}

	return path, nil
}

// SearchMenus searches menus by name or route
func (s *service) SearchMenus(ctx context.Context, query string) ([]MenuResponse, error) {
	if query == "" {
		return []MenuResponse{}, nil
	}

	// Get all menus and filter in memory (for simple implementation)
	// In production, you might want to add a search method to repository
	allMenus, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search menus: %w", err)
	}

	var results []MenuResponse
	searchTerm := strings.ToLower(strings.TrimSpace(query))

	s.searchInMenuTree(allMenus, searchTerm, &results)

	return results, nil
}

// CreateMenu creates a new menu with business validation
func (s *service) CreateMenu(ctx context.Context, input CreateMenuInput) (*MenuResponse, error) {
	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Check for duplicate name at the same level
	parentID := uint(0)
	if input.ParentID != nil {
		parentID = *input.ParentID
	}

	if err := s.checkDuplicateName(ctx, input.Name, parentID, 0); err != nil {
		return nil, err
	}

	// Check depth limit if parent is specified
	if input.ParentID != nil && *input.ParentID > 0 {
		depth, err := s.getMenuDepth(ctx, *input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check menu depth: %w", err)
		}
		if depth >= MaxMenuDepth {
			return nil, ErrMaxDepthExceeded
		}
	}

	// Generate menu code if not provided
	if input.Code == nil || *input.Code == "" {
		code := s.generateMenuCode(input.Name)
		input.Code = &code
	}

	// Create menu
	menu, err := s.repo.CreateMenu(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return menu, nil
}

// UpdateMenu updates an existing menu with validation
func (s *service) UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) (*MenuResponse, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	// Validate input
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Check for duplicate name if name is being updated
	if input.Name != nil {
		currentMenu, err := s.repo.GetMenuByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get current menu: %w", err)
		}

		parentID := currentMenu.ParentID
		if input.ParentID != nil {
			parentID = *input.ParentID
		}

		if err := s.checkDuplicateName(ctx, *input.Name, parentID, id); err != nil {
			return nil, err
		}
	}

	// Check depth limit if parent is being changed
	if input.ParentID != nil && *input.ParentID > 0 {
		depth, err := s.getMenuDepth(ctx, *input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check menu depth: %w", err)
		}
		if depth >= MaxMenuDepth {
			return nil, ErrMaxDepthExceeded
		}
	}

	// Update menu
	menu, err := s.repo.UpdateMenu(ctx, id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return menu, nil
}

// DeleteMenu deletes a menu with business validation
func (s *service) DeleteMenu(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}

	// Additional business logic before deletion
	// e.g., check if menu is being used in user permissions, etc.

	if err := s.repo.DeleteMenu(ctx, id); err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	return nil
}

// ReorderMenus updates the sort order of multiple menus
func (s *service) ReorderMenus(ctx context.Context, orders []MenuOrderInput) error {
	if len(orders) == 0 {
		return nil
	}

	// Validate all inputs first
	for _, order := range orders {
		if order.ID == 0 {
			return ErrInvalidInput
		}
		if order.Sort < MinSortValue || order.Sort > MaxSortValue {
			return ErrInvalidSortOrder
		}
	}

	// Update each menu's sort order
	for _, order := range orders {
		updateInput := UpdateMenuInput{
			Sort: &order.Sort,
		}

		_, err := s.repo.UpdateMenu(ctx, order.ID, updateInput)
		if err != nil {
			return fmt.Errorf("failed to reorder menu %d: %w", order.ID, err)
		}
	}

	return nil
}

// BulkUpdateMenuStatus updates the active status of multiple menus
func (s *service) BulkUpdateMenuStatus(ctx context.Context, ids []uint, isActive bool) error {
	if len(ids) == 0 {
		return nil
	}

	// Update each menu's status
	for _, id := range ids {
		if id == 0 {
			continue
		}

		updateInput := UpdateMenuInput{
			IsActive: &isActive,
		}

		_, err := s.repo.UpdateMenu(ctx, id, updateInput)
		if err != nil {
			return fmt.Errorf("failed to update menu %d status: %w", id, err)
		}
	}

	return nil
}

// ValidateMenuAccess checks if user has access to a menu (placeholder)
func (s *service) ValidateMenuAccess(ctx context.Context, userID uint, menuID uint) (bool, error) {
	if userID == 0 || menuID == 0 {
		return false, ErrInvalidInput
	}

	// Placeholder for permission checking logic
	// This would typically involve checking user roles, permissions, etc.
	// For now, return true (implement based on your permission system)

	return true, nil
}

// GetMenuStatistics returns menu usage statistics
func (s *service) GetMenuStatistics(ctx context.Context) (*MenuStatistics, error) {
	allMenus, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu statistics: %w", err)
	}

	stats := &MenuStatistics{
		TotalMenus:    0,
		ActiveMenus:   0,
		InactiveMenus: 0,
		RootMenus:     len(allMenus),
		MaxDepth:      0,
		GeneratedAt:   time.Now(),
	}

	s.calculateStatistics(allMenus, stats, 1)

	return stats, nil
}

// Private helper methods

// validateCreateInput validates create menu input
func (s *service) validateCreateInput(input CreateMenuInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	if len(input.Name) > MaxMenuNameLen {
		return fmt.Errorf("%w: name too long (max %d characters)", ErrInvalidInput, MaxMenuNameLen)
	}

	if input.Route != nil && len(*input.Route) > MaxRouteLen {
		return fmt.Errorf("%w: route too long (max %d characters)", ErrInvalidInput, MaxRouteLen)
	}

	if input.Sort < MinSortValue || input.Sort > MaxSortValue {
		return fmt.Errorf("%w: sort value must be between %d and %d", ErrInvalidInput, MinSortValue, MaxSortValue)
	}

	return nil
}

// validateUpdateInput validates update menu input
func (s *service) validateUpdateInput(input UpdateMenuInput) error {
	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
		}
		if len(*input.Name) > MaxMenuNameLen {
			return fmt.Errorf("%w: name too long (max %d characters)", ErrInvalidInput, MaxMenuNameLen)
		}
	}

	if input.Route != nil && len(*input.Route) > MaxRouteLen {
		return fmt.Errorf("%w: route too long (max %d characters)", ErrInvalidInput, MaxRouteLen)
	}

	if input.Sort != nil && (*input.Sort < MinSortValue || *input.Sort > MaxSortValue) {
		return fmt.Errorf("%w: sort value must be between %d and %d", ErrInvalidInput, MinSortValue, MaxSortValue)
	}

	return nil
}

// checkDuplicateName checks for duplicate menu names at the same level
func (s *service) checkDuplicateName(ctx context.Context, name string, parentID uint, excludeID uint) error {
	siblings, err := s.repo.GetMenusByParentID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("failed to check duplicate name: %w", err)
	}

	for _, sibling := range siblings {
		if sibling.ID != excludeID && strings.EqualFold(sibling.Name, name) {
			return ErrDuplicateMenuName
		}
	}

	return nil
}

// getMenuDepth calculates the depth of a menu in the hierarchy
func (s *service) getMenuDepth(ctx context.Context, menuID uint) (int, error) {
	depth := 0
	currentID := menuID

	for currentID != 0 {
		menu, err := s.repo.GetMenuByID(ctx, currentID)
		if err != nil {
			if errors.Is(err, ErrMenuNotFound) {
				break
			}
			return 0, err
		}
		depth++
		currentID = menu.ParentID

		// Prevent infinite loops
		if depth > MaxMenuDepth*2 {
			return 0, ErrCircularReference
		}
	}

	return depth, nil
}

// generateMenuCode generates a menu code from name
func (s *service) generateMenuCode(name string) string {
	// Simple code generation: lowercase, replace spaces with underscores
	code := strings.ToLower(strings.TrimSpace(name))
	code = strings.ReplaceAll(code, " ", "_")
	code = strings.ReplaceAll(code, "-", "_")

	// Remove special characters (keep only alphanumeric and underscores)
	var result strings.Builder
	for _, r := range code {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// enrichMenuData adds additional business data to menus
func (s *service) enrichMenuData(menus []MenuResponse) {
	// Add any additional business logic here
	// e.g., add user permissions, access flags, etc.
	for i := range menus {
		s.enrichSingleMenu(&menus[i])
		if len(menus[i].Children) > 0 {
			s.enrichMenuData(menus[i].Children)
		}
	}
}

// enrichSingleMenu enriches a single menu item
func (s *service) enrichSingleMenu(menu *MenuResponse) {
	// Add business logic here
	// e.g., set access flags, calculate permissions, etc.
}

// searchInMenuTree recursively searches through menu tree
func (s *service) searchInMenuTree(menus []MenuResponse, searchTerm string, results *[]MenuResponse) {
	for _, menu := range menus {
		// Check if menu matches search criteria
		if s.menuMatchesSearch(menu, searchTerm) {
			*results = append(*results, menu)
		}

		// Search in children
		if len(menu.Children) > 0 {
			s.searchInMenuTree(menu.Children, searchTerm, results)
		}
	}
}

// menuMatchesSearch checks if a menu matches the search criteria
func (s *service) menuMatchesSearch(menu MenuResponse, searchTerm string) bool {
	nameMatch := strings.Contains(strings.ToLower(menu.Name), searchTerm)

	routeMatch := false
	if menu.Route != nil {
		routeMatch = strings.Contains(strings.ToLower(*menu.Route), searchTerm)
	}

	return nameMatch || routeMatch
}

// calculateStatistics recursively calculates menu statistics
func (s *service) calculateStatistics(menus []MenuResponse, stats *MenuStatistics, currentDepth int) {
	for _, menu := range menus {
		stats.TotalMenus++

		if menu.IsActive {
			stats.ActiveMenus++
		} else {
			stats.InactiveMenus++
		}

		if currentDepth > stats.MaxDepth {
			stats.MaxDepth = currentDepth
		}

		if len(menu.Children) > 0 {
			s.calculateStatistics(menu.Children, stats, currentDepth+1)
		}
	}
}
