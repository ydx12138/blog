package dto

import (
	"blog/models"
	"time"
)

type ArticleSimple struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	Title        string `gorm:"size:200;not null;comment:标题" json:"title"`
	Summary      string `gorm:"size:500;comment:摘要" json:"summary"`
	Cover        string `gorm:"size:255;comment:封面" json:"cover"`
	CategoryID   uint64 `json:"-" gorm:"comment:类别ID" json:"category_id"`
	Category     models.Category
	ViewCount    uint64 `gorm:"default:0;comment:浏览数" json:"view_count"`
	LikeCount    uint64 `gorm:"default:0;comment:点赞数" json:"like_count"`
	CommentCount uint64 `gorm:"default:0;comment:评论数" json:"comment_count"`

	Tags string `gorm:"size:255;comment:标签，逗号分隔" json:"tags"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
