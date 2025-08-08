package model

import (
	"time"

	"gorm.io/gorm"
)

type People struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PersonID  string         `json:"person_id" gorm:"column:person_id;type:varchar(150);not null"`
	Username  string         `json:"username" gorm:"column:username;type:varchar(150);not null"`
	Password  string         `json:"password" gorm:"column:password;type:varchar(150);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (People) TableName() string {
	return "peoples"
}
