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
	err := core.DB.Model(models.Article{}).Preload("Category").Limit(10).Offset((page-1)*10).Find(&articleList, "status = ?", 2).Error
	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		return nil, err
	}
	return articleList, nil
}
