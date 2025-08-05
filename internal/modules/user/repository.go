package user

import (
	"kswi-backend/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
