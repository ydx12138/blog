package admin

import (
	"blog/internal/dao"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllComments 获取全部评论
func GetAllComments(c *gin.Context) {
	var q dto.AdminCommentQuery
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
	keyword := c.Query("keyword")
	searchType := c.Query("type") // "nickname" or "content"
	comments, total, err := dao.GetAllComments(q.Page, q.PageSize, keyword, searchType)
	if err != nil {
		zap.L().Error("GetAllComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  comments,
		"total": total,
	}, c)
}

// GetPendingComments 获取待审核评论
func GetPendingComments(c *gin.Context) {
	var q dto.AdminCommentQuery
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
	comments, total, err := dao.GetPendingComments(q.Page, q.PageSize)
	if err != nil {
		zap.L().Error("GetPendingComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  comments,
		"total": total,
	}, c)
}

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

// ApproveComment 通过评论
func ApproveComment(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 获取评论信息
	comment, err := dao.GetCommentByID(id)
	if err != nil {
		zap.L().Error("ApproveComment 评论不存在:" + err.Error())
		response.ErrWithMsg(code.ErrCommentNotFound, c)
		return
	}
	// 只有待审核(3)的评论通过才 +1
	if comment.Status == 3 {
		_ = dao.UpdateArticleCommentCount(comment.ArticleID, 1)
	}
	err = dao.UpdateCommentStatus(id, 1)
	if err != nil {
		zap.L().Error("ApproveComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("审核通过", c)
}

// RejectComment 拒绝评论
func RejectComment(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.UpdateCommentStatus(id, 2)
	if err != nil {
		zap.L().Error("RejectComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("已拒绝", c)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 获取评论信息，已通过的评论删除时 -1
	comment, err := dao.GetCommentByID(id)
	if err == nil && comment.Status == 1 {
		_ = dao.UpdateArticleCommentCount(comment.ArticleID, -1)
	}
	err = dao.DeleteComment(id)
	if err != nil {
		zap.L().Error("DeleteComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}

// SetCommentStatus 设置评论状态（1正常 2隐藏 3待审核）
func SetCommentStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if c.ShouldBind(&body) != nil || (body.Status < 1 || body.Status > 3) {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.UpdateCommentStatus(id, body.Status)
	if err != nil {
		zap.L().Error("SetCommentStatus:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("状态已更新", c)
}
