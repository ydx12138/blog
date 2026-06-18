package dao

import (
	"blog/core"
	"blog/models"
	"blog/models/dto"

	"go.uber.org/zap"
)

// 按页码获取10条已发布的文章
func GetArticleByPage(page int) ([]dto.ArticleSimple, error) {
	var articleList []dto.ArticleSimple = make([]dto.ArticleSimple, 0)
	//err := core.DB.Model(models.Article{}).Preload("Category").Limit(10).Offset((page-1)*10).Find(&articleList, "status = ?", 2).Error
	err := core.DB.
		Model(models.Article{}).
		Select(`
			article.id,
			article.title,
			article.summary,
			article.cover,
			article.category_id,
			c.name AS category_name,
			article.view_count,
			article.like_count,
			article.comment_count,
			article.tags,
			article.created_at,
			article.updated_at
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Limit(10).
		Offset((page - 1) * 10).
		Scan(&articleList).Error

	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		return articleList, err
	}
	return articleList, nil
}

func GetArticleDetail(id uint64) (dto.ArticleDetail, error) {
	//core.DB.Model(models.Article{}).Find(&detail)
	var detail dto.ArticleDetail
	err := core.DB.
		Model(models.Article{}).
		Select(`
			article.id,
			article.title,
			article.summary,
			article.content,
			article.cover,
			c.name AS category_name,
			article.view_count,
			article.like_count,
			article.comment_count,
			article.publish_time,
			article.tags
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Where("article.id = ?", id).
		Scan(&detail).Error
	if err != nil {
		zap.L().Error("GetArticleDetail()" + err.Error())
		return detail, err
	}
	return detail, nil
}

func SearchArticleByKey(keyword string) ([]dto.ArticleSimple, error) {
	var articleList []dto.ArticleSimple = make([]dto.ArticleSimple, 0)
	err := core.DB.
		Model(models.Article{}).
		Select(`
			article.*,
			c.name AS category_name
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Where("article.title like ? or article.summary like ? or article.content like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&articleList).Error
	if err != nil {
		zap.L().Error("SearchArticleByKey()" + err.Error())
		return articleList, err
	}
	return articleList, nil
}
