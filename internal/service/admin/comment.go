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
	err = dao.DeleteComment(id)
	if err != nil {
		zap.L().Error("DeleteComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}
