package oss

import "context"

type Service interface {
	DtDatabase(ctx context.Context, req DtDatabaseRequest) ([]DtDatabaseResponse, int, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) DtDatabase(ctx context.Context, req DtDatabaseRequest) ([]DtDatabaseResponse, int, int, error) {
	return s.repo.DtDatabase(ctx, req)
}
