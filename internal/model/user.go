package model

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	PersonID  string    `json:"person_id" gorm:"column:person_id;not null"`
	Username  string    `json:"username" gorm:"column:username;not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password;not null"` // Use '-' to exclude from JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
