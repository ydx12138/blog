package service

import (
	"blog/internal/dao"
	"blog/models/request"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 首页返回文章列表
func GetArticles(c *gin.Context) {
	//取参
	var p request.PageQuery
	err := c.ShouldBindQuery(&p)
	//参数出错
	if err != nil {
		zap.L().Error("GetArticles分页参数错误" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//sql
	articleSimpleList, err := dao.GetArticleByPage(p.Page)
	//sql出错
	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	//返回文章列表
	response.SuccessWithData(articleSimpleList, c)
}
