package models

import "time"

type Comment struct {
	ID uint64 `gorm:"primaryKey"`

	ArticleID uint64 `gorm:"comment:文章ID"`
	Article   Article

	UserID uint64 `gorm:"comment:评论所属用户ID"`
	User   User

	ParentID uint64 `gorm:"default:0;comment:父评论ID"`

	Content string `gorm:"type:text;not null;comment:评论内容"`

	Status int8 `gorm:"default:1;comment:1正常 2隐藏 3待审核"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
