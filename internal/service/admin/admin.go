package admin

import (
	"blog/internal/dao"
	"blog/internal/utils"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 管理员登录
func Login(c *gin.Context) {
	var ad dto.AdminLogin
	err := c.ShouldBind(&ad)
	if err != nil {
		zap.L().Error("Admin Login 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	result, err := dao.LoginVerification(ad.Username, ad.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrUserNotFound, c)
		} else if err.Error() == "密码错误" {
			response.ErrWithMsg(code.ErrPassword, c)
		} else {
			zap.L().Error("Admin Login:" + err.Error())
			response.ErrWithMsg(code.InternalError, c)
		}
		return
	}
	// 生成token
	token, err := utils.GenerateToken(result.ID, "admin")
	if err != nil {
		zap.L().Error("Admin Login 生成Token失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}

	response.SuccessWithData(map[string]interface{}{
		"token":    token,
		"nickname": result.Nickname,
		"username": result.Username,
	}, c)
}
