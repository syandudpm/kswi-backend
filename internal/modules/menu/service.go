package menu

import (
	"context"
)

type Service interface {
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	CreateMenu(ctx context.Context, input CreateMenuInput) error
	UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) error
	DeleteMenu(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	return s.repo.GetMenuTree(ctx)
}

func (s *service) CreateMenu(ctx context.Context, input CreateMenuInput) error {
	return s.repo.CreateMenu(ctx, input)
}

func (s *service) UpdateMenu(ctx context.Context, id uint, input UpdateMenuInput) error {
	return s.repo.UpdateMenu(ctx, id, input)
}

func (s *service) DeleteMenu(ctx context.Context, id uint) error {
	return s.repo.DeleteMenu(ctx, id)
}
