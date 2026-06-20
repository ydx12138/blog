package models

import "time"

type Article struct {
	ID      uint64 `gorm:"primaryKey"`
	Title   string `gorm:"size:200;not null;comment:标题"`
	Summary string `gorm:"size:500;comment:摘要"`

	Content     string `gorm:"type:longtext;comment:内容"`
	ContentType int8   `gorm:"default:1;comment:1富文本 2markdown"`

	Cover string `gorm:"size:255;comment:封面"`

	CategoryID uint64   `gorm:"comment:类别ID"`
	Category   Category //不在表里

	AuthorID uint64 `gorm:"comment:作者ID"`
	Author   Admin  //不在表里

	ViewCount    uint64 `gorm:"default:0;comment:浏览数"`
	LikeCount    uint64 `gorm:"default:0;comment:点赞数"`
	CommentCount uint64 `gorm:"default:0;comment:评论数"`

	Status int8 `gorm:"default:1;comment:1草稿 2发布"`

	PublishTime *time.Time

	Tags string `gorm:"size:255;comment:标签，逗号分隔"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
