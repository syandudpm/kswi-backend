package person

import (
	"kswi-backend/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(person *model.Person) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(person *model.Person) error {
	return r.db.Create(person).Error
}
