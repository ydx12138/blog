package admin

import (
	"blog/internal/dao"
	"blog/models"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetArticles 管理端获取文章列表
func GetArticles(c *gin.Context) {
	var q dto.AdminArticleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		q.Page = 1
		q.PageSize = 10
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	articles, total, err := dao.AdminGetArticles(q.Page, q.PageSize, q.Keyword, q.Status)
	if err != nil {
		zap.L().Error("Admin GetArticles:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  articles,
		"total": total,
	}, c)
}

// GetArticle 管理端获取单篇文章
func GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		zap.L().Error("Admin GetArticle:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	article, err := dao.GetArticleByID(id)
	if err != nil {
		zap.L().Error("Admin GetArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(article, c)
}

// CreateArticle 创建文章
func CreateArticle(c *gin.Context) {
	var req dto.CreateArticleReq
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("CreateArticle 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 从JWT获取作者ID
	authorID, _ := c.Get("userID")
	// 未选分类默认归入"杂谈"
	categoryID := req.CategoryID
	if categoryID == 0 {
		defaultCat, err := dao.GetOrCreateDefaultCategory()
		if err == nil {
			categoryID = defaultCat.ID
		}
	}
	article := models.Article{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		ContentType: req.ContentType,
		Cover:       req.Cover,
		CategoryID:  categoryID,
		Tags:        req.Tags,
		Status:      req.Status,
		AuthorID:    authorID.(uint64),
	}
	//如果状态是发布，直接就把发布时间设置为当前时间
	if req.Status == 2 {
		now := time.Now()
		article.PublishTime = &now
	}
	err = dao.CreateArticle(&article)
	if err != nil {
		zap.L().Error("CreateArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("创建成功", c)
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	var req dto.UpdateArticleReq
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("UpdateArticle 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 获取现有文章
	article, err := dao.GetArticleByID(req.ID)
	if err != nil {
		zap.L().Error("UpdateArticle 文章不存在:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	// 更新字段
	article.Title = req.Title
	article.Summary = req.Summary
	article.Content = req.Content
	article.ContentType = req.ContentType
	article.Cover = req.Cover
	article.CategoryID = req.CategoryID
	article.Tags = req.Tags
	article.Status = req.Status
	if req.Status == 2 && article.PublishTime == nil {
		now := time.Now()
		article.PublishTime = &now
	}
	err = dao.UpdateArticle(&article)
	if err != nil {
		zap.L().Error("UpdateArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("更新成功", c)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	var req dto.IDReq
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("DeleteArticle:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.DeleteArticle(req.ID)
	if err != nil {
		zap.L().Error("DeleteArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}

// GetDrafts 获取草稿列表
func GetDrafts(c *gin.Context) {
	var q dto.AdminArticleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		q.Page = 1
		q.PageSize = 10
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	articles, total, err := dao.GetDrafts(q.Page, q.PageSize)
	if err != nil {
		zap.L().Error("GetDrafts:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  articles,
		"total": total,
	}, c)
}

// PublishArticle 发布文章
func PublishArticle(c *gin.Context) {
	var req dto.IDReq
	err := c.ShouldBind(&req)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	article, err := dao.GetArticleByID(req.ID)
	if err != nil {
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	article.Status = 2
	now := time.Now()
	article.PublishTime = &now
	err = dao.UpdateArticle(&article)
	if err != nil {
		zap.L().Error("PublishArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("发布成功", c)
}
