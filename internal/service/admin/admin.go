package admin

import (
	"blog/internal/dao"
	"blog/internal/utils"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 管理员登录
func Login(c *gin.Context) {
	//参数
	var ad dto.AdminLogin
	err := c.ShouldBind(&ad)
	if err != nil {
		zap.L().Error("Login:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//sql
	result, err := dao.LoginVerification(ad.Username, ad.Password)
	if err != nil {
		zap.L().Error("Login:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	if result.ID == 0 {
		zap.L().Info("管理员登录失败")
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	//token
	token, err := utils.GenerateToken(result.ID, "admin")
	if err != nil {
		zap.L().Error("Login:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}

	response.SuccessWithData(map[string]string{"token": token}, c)
}
