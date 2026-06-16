package models

import "time"

//type User struct {
//	ID       uint64 `gorm:"primaryKey"`
//	Username string `gorm:"size:50;uniqueIndex;not null"`
//	Password string `gorm:"size:255;not null"`
//	Nickname string `gorm:"size:50;not null"`
//	Avatar   string `gorm:"size:255;default:''"`
//	Email    string `gorm:"size:100;uniqueIndex;not null"`
//	Status   int8   `gorm:"default:1;comment:1正常 2禁用"`
//
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Username  string `gorm:"size:50;uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	Nickname  string `gorm:"size:50"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt time.Time
}
