package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// User fields
	Username  string `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string `json:"-" gorm:"not null;size:255"`
	FirstName string `json:"first_name" gorm:"size:50"`
	LastName  string `json:"last_name" gorm:"size:50"`
	Role      string `json:"role" gorm:"not null;default:'user';size:20"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`

	// Profile fields
	Avatar      string     `json:"avatar,omitempty" gorm:"size:255"`
	Phone       string     `json:"phone,omitempty" gorm:"size:20"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`

	// Metadata
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count" gorm:"default:0"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}
