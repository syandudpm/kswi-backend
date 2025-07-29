package menu

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) MenuTree(ctx context.Context) ([]MenuResponse, error) {
	return s.repo.GetMenuTree(ctx)
}
