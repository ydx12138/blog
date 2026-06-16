package models

import "time"

type Admin struct {
	ID       uint64 `gorm:"primaryKey;"`
	Username string `gorm:"size:50;uniqueIndex;not null;comment:用户名"`
	Password string `gorm:"size:255;not null;comment:密码"`
	Nickname string `gorm:"size:50;comment:昵称"`
	Email    string `gorm:"size:100;comment:邮箱"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
