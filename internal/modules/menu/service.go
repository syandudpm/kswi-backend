package menu

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// MenuTree returns the complete menu tree as DTOs
func (s *Service) MenuTree(ctx context.Context) ([]MenuResponse, error) {
	// Get menu tree from repository (returns models)
	menuTree, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu tree: %w", err)
	}

	// Convert models to DTOs
	return ToMenuResponseList(menuTree), nil
}

// GetActiveMenus returns all active menus (flat list)
func (s *Service) GetActiveMenus(ctx context.Context) ([]MenuResponse, error) {
	menus, err := s.repo.GetActiveMenus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active menus: %w", err)
	}

	return ToMenuResponseList(menus), nil
}

// GetMenuByID returns a specific menu by ID
func (s *Service) GetMenuByID(ctx context.Context, id uint) (*MenuResponse, error) {
	menu, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu by ID: %w", err)
	}

	response := ToMenuResponse(*menu)
	return &response, nil
}
