package models

import "time"

type Category struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"size:50;not null"`
	Description string `gorm:"size:255"`
	Sort        int    `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
