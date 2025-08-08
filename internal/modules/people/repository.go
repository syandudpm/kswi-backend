package people

import (
	"kswi-backend/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(people *model.People) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(people *model.People) error {
	return r.db.Create(people).Error
}
