package user

import (
	"blog/internal/dao"
	"blog/models"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetComments 获取文章评论列表（公开）
func GetComments(c *gin.Context) {
	var q dto.CommentListQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		zap.L().Error("GetComments 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if q.Page < 1 {
		q.Page = 1
	}
	comments, total, err := dao.GetCommentsByArticle(q.ArticleID, q.Page, 10)
	if err != nil {
		zap.L().Error("GetComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  comments,
		"total": total,
	}, c)
}

// CreateComment 创建评论（需要登录）
func CreateComment(c *gin.Context) {
	var req dto.CreateCommentReq
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("CreateComment 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 从JWT获取userID
	userID, exists := c.Get("userID")
	if !exists {
		response.ErrWithMsg(code.Unauthorized, c)
		return
	}
	comment := models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID.(uint64),
		Content:   req.Content,
		ParentID:  req.ParentID,
		Status:    3, // 默认待审核
	}
	err = dao.CreateComment(&comment)
	if err != nil {
		zap.L().Error("CreateComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("评论成功，审核通过后显示", c)
}
