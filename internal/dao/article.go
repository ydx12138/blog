package dao

import (
	"blog/core"
	"blog/models/dto"

	"go.uber.org/zap"
)

// 按页码获取10条已发布的文章
func GetArticleByPage(page int) ([]dto.ArticleSimple, error) {
	var articleList []dto.ArticleSimple = make([]dto.ArticleSimple, 0)
	//err := core.DB.Model(models.Article{}).Preload("Category").Limit(10).Offset((page-1)*10).Find(&articleList, "status = ?", 2).Error
	err := core.DB.
		Table("article a").
		Select(`
			a.id,
			a.title,
			a.summary,
			a.cover,
			a.category_id,
			c.name AS category_name,
			a.view_count,
			a.like_count,
			a.comment_count,
			a.tags,
			a.created_at,
			a.updated_at
		`).
		Joins("LEFT JOIN category c ON a.category_id = c.id").
		Where("a.status = ?", 2).
		Limit(10).
		Offset((page - 1) * 10).
		Scan(&articleList).Error

	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		return nil, err
	}
	return articleList, nil
}
