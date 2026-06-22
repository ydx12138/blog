package dao

import (
	"blog/core"
	"blog/models"
	"blog/models/vo"
	"strings"

	"go.uber.org/zap"
)

// 按页码获取已发布的文章(支持自定义pageSize)
func GetArticleByPage(page int, pageSize int) ([]vo.ArticleSimple, error) {
	var articleList []vo.ArticleSimple = make([]vo.ArticleSimple, 0)
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
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&articleList).Error

	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		return articleList, err
	}
	return articleList, nil
}

func GetArticleDetail(id uint64) (vo.ArticleDetail, error) {
	var detail vo.ArticleDetail
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
			article.tags,
			article.content_type
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

func SearchArticleByKey(keyword string) ([]vo.ArticleSimple, error) {
	var articleList []vo.ArticleSimple = make([]vo.ArticleSimple, 0)
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
		Where("article.title like ? or article.summary like ? or article.content like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&articleList).Error
	if err != nil {
		zap.L().Error("SearchArticleByKey()" + err.Error())
		return articleList, err
	}
	return articleList, nil
}

// GetArticleByCategory 按分类获取已发布文章
func GetArticleByCategory(categoryID uint64, page int, pageSize int) ([]vo.ArticleSimple, error) {
	var articleList []vo.ArticleSimple = make([]vo.ArticleSimple, 0)
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
		Where("article.category_id = ?", categoryID).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&articleList).Error
	if err != nil {
		zap.L().Error("GetArticleByCategory()" + err.Error())
		return articleList, err
	}
	return articleList, nil
}

// IncrementViewCount 增加文章浏览数
func IncrementViewCount(id uint64) error {
	err := core.DB.Model(&models.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", core.DB.Raw("view_count + ?", 1)).Error
	if err != nil {
		zap.L().Error("IncrementViewCount()" + err.Error())
	}
	return err
}

// ========== 管理端 Article DAO ==========

// AdminGetArticles 管理端获取文章列表（含搜索和状态筛选）
func AdminGetArticles(page int, pageSize int, keyword string, status int8) ([]models.Article, int64, error) {
	var articles []models.Article = make([]models.Article, 0)
	var total int64
	query := core.DB.Model(&models.Article{}).Preload("Category")
	if keyword != "" {
		query = query.Where("title like ? or summary like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("AdminGetArticles count:" + err.Error())
		return articles, total, err
	}
	err = query.Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&articles).Error
	if err != nil {
		zap.L().Error("AdminGetArticles:" + err.Error())
		return articles, total, err
	}
	return articles, total, nil
}

// GetArticleByID 根据ID获取完整文章模型（供管理端编辑使用）
func GetArticleByID(id uint64) (models.Article, error) {
	var article models.Article
	err := core.DB.Preload("Category").First(&article, id).Error
	if err != nil {
		zap.L().Error("GetArticleByID:" + err.Error())
		return article, err
	}
	return article, nil
}

// CreateArticle 创建文章
func CreateArticle(article *models.Article) error {
	err := core.DB.Create(article).Error
	if err != nil {
		zap.L().Error("CreateArticle:" + err.Error())
		return err
	}
	return nil
}

// UpdateArticle 更新文章
func UpdateArticle(article *models.Article) error {
	err := core.DB.Save(article).Error
	if err != nil {
		zap.L().Error("UpdateArticle:" + err.Error())
		return err
	}
	return nil
}

// DeleteArticle 删除文章
func DeleteArticle(id uint64) error {
	err := core.DB.Delete(&models.Article{}, id).Error
	if err != nil {
		zap.L().Error("DeleteArticle:" + err.Error())
		return err
	}
	return nil
}

// GetDrafts 获取草稿列表
func GetDrafts(page int, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article = make([]models.Article, 0)
	var total int64
	err := core.DB.Model(&models.Article{}).Preload("Category").
		Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return articles, total, err
	}
	err = core.DB.Preload("Category").Where("status = ?", 1).
		Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&articles).Error
	if err != nil {
		zap.L().Error("GetDrafts:" + err.Error())
		return articles, total, err
	}
	return articles, total, nil
}

// GetAllTags 获取所有已使用的标签（去重）
func GetAllTags() ([]string, error) {
	var articles []models.Article
	err := core.DB.Select("tags").Where("status = ? AND tags != ''", 2).Find(&articles).Error
	if err != nil {
		zap.L().Error("GetAllTags:" + err.Error())
		return nil, err
	}
	tagSet := make(map[string]struct{})
	for _, a := range articles {
		for _, t := range splitTags(a.Tags) {
			tagSet[t] = struct{}{}
		}
	}
	result := make([]string, 0, len(tagSet))
	for t := range tagSet {
		result = append(result, t)
	}
	return result, nil
}

func splitTags(tags string) []string {
	parts := make([]string, 0)
	for _, t := range strings.Split(tags, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			parts = append(parts, t)
		}
	}
	return parts
}

// IncrementLikeCount 文章点赞数+1
func IncrementLikeCount(articleID uint64) error {
	err := core.DB.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", core.DB.Raw("like_count + ?", 1)).Error
	if err != nil {
		zap.L().Error("IncrementLikeCount:" + err.Error())
	}
	return err
}

// UpdateArticleCommentCount 更新文章评论数
func UpdateArticleCommentCount(articleID uint64, delta int) error {
	err := core.DB.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("comment_count", core.DB.Raw("comment_count + ?", delta)).Error
	if err != nil {
		zap.L().Error("UpdateArticleCommentCount:" + err.Error())
	}
	return err
}
