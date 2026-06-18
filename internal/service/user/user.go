package user

import (
	"blog/internal/dao"
	"blog/models/dto"
	"blog/models/vo"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 首页返回按页码返回文章列表
func GetArticles(c *gin.Context) {
	//取参
	var p dto.PageQuery
	err := c.ShouldBindQuery(&p)
	//参数出错
	if err != nil {
		zap.L().Error("GetArticles参数错误" + err.Error())
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

// 根据id获取文章详情
func GetArticle(c *gin.Context) {
	//参数
	var detail vo.ArticleDetail
	err := c.ShouldBindQuery(&detail)
	if err != nil {
		zap.L().Error("GetArticlec参数错误" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//sql
	articleDetail, err := dao.GetArticleDetail(detail.ID)
	//结果
	if err != nil {
		zap.L().Error("GetArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articleDetail, c)
}

// 文章搜索
func SearchArticle(c *gin.Context) {
	var key dto.ArticleKeyWord
	err := c.ShouldBindQuery(&key)
	if err != nil {
		zap.L().Error("PostArticle参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//sql
	articleSimples, err := dao.SearchArticleByKey(key.Keyword)
	if err != nil {
		zap.L().Error("PostArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articleSimples, c)
}
