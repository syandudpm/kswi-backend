package menu

import (
	"context"
	"errors"
	"sort"

	"gorm.io/gorm"
)

type Repository interface {
	GetMenuTree(ctx context.Context) ([]MenuResponse, error)
	FindByID(ctx context.Context, id uint) (*MenuDetailResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetMenuTree(ctx context.Context) ([]MenuResponse, error) {
	var menus []MenuResponse

	err := r.db.WithContext(ctx).
		Table("menus").
		Select("id, parent_id, sort, name, route, icon").
		Where("is_active = ?", true).
		Order("sort ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	return r.buildMenuTree(menus, 0), nil
}

func (r *repository) buildMenuTree(flat []MenuResponse, parentID uint) []MenuResponse {
	var result []MenuResponse

	for _, item := range flat {
		if item.ParentID == parentID {
			children := r.buildMenuTree(flat, item.ID)

			if children == nil {
				children = []MenuResponse{}
			}

			item.Children = children
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Sort < result[j].Sort
	})

	return result
}

func (r *repository) FindByID(ctx context.Context, id uint) (*MenuDetailResponse, error) {
	var dto MenuDetailResponse

	err := r.db.WithContext(ctx).
		Table("menus").
		Select("id, parent_id, sort, name, route, icon, is_active").
		Where("id = ? AND is_active = ?", id, true).
		First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &dto, nil
}
