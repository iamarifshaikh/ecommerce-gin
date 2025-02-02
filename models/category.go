package models

import "time"

type Category struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null;unique"`
	ParentID      *uint
	Subcategories []Category `gorm:"foreignKey:ParentID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
