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

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
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
	keyword := c.Query("keyword")
	statusStr := c.Query("status")
	var status uint64
	if statusStr != "" {
		status, _ = strconv.ParseUint(statusStr, 10, 64)
	}
	users, total, err := dao.GetUsersByPage(q.Page, q.PageSize, keyword, status)
	if err != nil {
		zap.L().Error("GetUsers:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	type SafeUser struct {
		ID        uint64 `json:"id"`
		Email     string `json:"email"`
		Nickname  string `json:"nickname"`
		Status    uint64 `json:"status"`
		CreatedAt string `json:"created_at"`
	}
	var safeUsers []SafeUser = make([]SafeUser, 0)
	for _, u := range users {
		safeUsers = append(safeUsers, SafeUser{
			ID:        u.ID,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.SuccessWithData(map[string]interface{}{
		"list":  safeUsers,
		"total": total,
	}, c)
}

func parseAdminID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	id, err := parseAdminID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.UpdateUserStatus(id, 2)
	if err != nil {
		zap.L().Error("BanUser:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("封禁成功", c)
}

// UnbanUser 解封用户
func UnbanUser(c *gin.Context) {
	id, err := parseAdminID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.UpdateUserStatus(id, 1)
	if err != nil {
		zap.L().Error("UnbanUser:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("解封成功", c)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := parseAdminID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.DeleteUserByID(id)
	if err != nil {
		zap.L().Error("DeleteUser:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}
