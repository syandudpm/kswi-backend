package model

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ParentID     uint    `gorm:"not null;default:0;index:idx_menus_parent_id"`
	PermissionID *uint   `gorm:"index:idx_menus_permission_id"`
	Code         *string `gorm:"size:200;index:idx_menus_code"`
	ParentCode   *string `gorm:"size:200"`
	Sort         int     `gorm:"not null;default:0"`
	Name         string  `gorm:"not null;size:200"`
	Route        *string `gorm:"size:200"`
	Icon         *string `gorm:"size:1000"`
	IsActive     bool    `gorm:"not null;default:true"`

	// Relations
	Parent   *Menu  `gorm:"foreignKey:ParentID"`
	Children []Menu `gorm:"foreignKey:ParentID"`
}

func (Menu) TableName() string {
	return "menus"
}
