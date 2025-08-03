package users

import (
	"context"
	"fmt"

	"kswi-backend/internal/config"
	"kswi-backend/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	// List(ctx context.Context, params pagination.Params) ([]model.User, *pagination.Meta, error)
	Count(ctx context.Context) (int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository() Repository {
	return &repository{
		db: config.GetDB(),
	}
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *repository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with pagination
// func (r *repository) List(ctx context.Context, params pagination.Params) ([]model.User, *pagination.Meta, error) {
// 	var users []model.User
// 	var total int64

// 	// Count total records
// 	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
// 		return nil, nil, fmt.Errorf("failed to count users: %w", err)
// 	}

// 	// Calculate pagination
// 	meta := pagination.Calculate(params, total)

// 	// Query with pagination
// 	query := r.db.WithContext(ctx).
// 		Offset(meta.Offset).
// 		Limit(meta.Limit)

// 	// Add sorting
// 	if params.Sort != "" {
// 		query = query.Order(params.Sort + " " + params.Order)
// 	} else {
// 		query = query.Order("created_at DESC")
// 	}

// 	if err := query.Find(&users).Error; err != nil {
// 		return nil, nil, fmt.Errorf("failed to list users: %w", err)
// 	}

// 	return users, meta, nil
// }

// Count returns the total number of users
func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
