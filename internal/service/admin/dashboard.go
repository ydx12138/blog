package admin

import (
	"blog/core"
	"blog/models"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardData struct {
	TotalArticles     int64 `json:"total_articles"`
	PublishedArticles int64 `json:"published_articles"`
	DraftArticles     int64 `json:"draft_articles"`
	TotalComments     int64 `json:"total_comments"`
	PendingComments   int64 `json:"pending_comments"`
	TotalUsers        int64 `json:"total_users"`
	TotalViews        int64 `json:"total_views"`
}

func GetDashboard(c *gin.Context) {
	var data DashboardData

	core.DB.Model(&models.Article{}).Count(&data.TotalArticles)
	core.DB.Model(&models.Article{}).Where("status = ?", 2).Count(&data.PublishedArticles)
	core.DB.Model(&models.Article{}).Where("status = ?", 1).Count(&data.DraftArticles)
	core.DB.Model(&models.Comment{}).Count(&data.TotalComments)
	core.DB.Model(&models.Comment{}).Where("status = ?", 3).Count(&data.PendingComments)
	core.DB.Model(&models.User{}).Count(&data.TotalUsers)

	var totalViews int64
	core.DB.Model(&models.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	data.TotalViews = totalViews

	if core.DB.Error != nil {
		zap.L().Error("GetDashboard error")
	}

	response.SuccessWithData(data, c)
}
