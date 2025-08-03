package menu

import (
	"context"
	"fmt"
	"kswi-backend/internal/shared/errors"
)

type Service interface {
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	GetMenuByID(ctx context.Context, id uint) (*MenuDetailResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetMenuTree retrieves the hierarchical menu structure
func (s *service) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	menus, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get menu tree: %w", err))
	}

	return menus, nil
}

func (s *service) GetMenuByID(ctx context.Context, id uint) (*MenuDetailResponse, error) {
	menu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo error: %w", err)
	}
	if menu == nil {
		return nil, nil // Not found
	}
	return menu, nil
}
