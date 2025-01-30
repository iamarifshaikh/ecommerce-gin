package models

import "time"

type BlacklistToken struct {
	ID        uint      `gorm:"primary_key"`
	Token     string    `gorm:"type:text; not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
