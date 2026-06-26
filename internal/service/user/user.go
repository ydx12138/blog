package user

import (
	"blog/internal/dao"
	"blog/internal/utils"
	"blog/models"
	"blog/models/dto"
	"blog/models/vo"
	"blog/pkg/code"
	"blog/pkg/response"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 首页返回按页码返回文章列表
func GetArticles(c *gin.Context) {
	//参数接收和校验处理
	var p dto.PageQueryWithSize
	err := c.ShouldBindQuery(&p)
	if err != nil {
		zap.L().Error("GetArticles参数page出错，已采用默认参数page=1:" + err.Error())
		p.Page = 1
	}
	if p.Page < 1 {
		zap.L().Error("GetArticles参数page出错，已采用默认参数page=1")
		p.Page = 1
	}
	if p.PageSize < 1 {
		zap.L().Error("GetArticles参数PageSize出错，已采用默认参数PageSize=10")
		p.PageSize = 10
	}
	//db
	articles, total, err := dao.GetArticleByPage(p.Page, p.PageSize)
	if err != nil {
		zap.L().Error("GetArticleByPage()" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": articles, "total": total}, c)
}

// 根据id获取文章详情
func GetArticle(c *gin.Context) {
	//接收、校验参数
	var detail vo.ArticleDetail
	err := c.ShouldBindQuery(&detail)
	if err != nil {
		zap.L().Error("GetArticle 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//db
	articleDetail, err := dao.GetArticleDetail(detail.ID)
	if err != nil {
		zap.L().Error("GetArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	// 增加浏览数
	_ = dao.IncrementViewCount(detail.ID)
	response.SuccessWithData(articleDetail, c)
}

// 文章搜索
func SearchArticle(c *gin.Context) {
	var key dto.ArticleKeyWord
	err := c.ShouldBindQuery(&key)
	if err != nil {
		zap.L().Error("SearchArticle 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	articleSimples, err := dao.SearchArticleByKey(key.Keyword)
	if err != nil {
		zap.L().Error("SearchArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articleSimples, c)
}

// 注册
func Register(c *gin.Context) {
	var user dto.UserRegister
	err := c.ShouldBind(&user)
	if err != nil {
		zap.L().Error("Register:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 判断邮箱是否已注册
	existUser, err := dao.GetUserByEmail(user.Email)
	//db出错，ErrRecordNotFound除外
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Error("Register 查询用户失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	//
	if existUser.ID != 0 {
		zap.L().Error("邮箱已经存在")
		response.ErrWithMsg(code.ErrUserExist, c)
		return
	}
	// 加密
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		zap.L().Error("Register 密码加密失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	// 存入数据库
	newUser := models.User{
		Email:    user.Email,
		Password: hashPassword,
		Nickname: user.Nickname,
	}
	err = dao.CreateUser(newUser)
	if err != nil {
		zap.L().Error("Register:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	fmt.Println("4", user)
	response.SuccessWithMsg("注册成功", c)
}

// 用户登录
func Login(c *gin.Context) {
	var login dto.UserLogin
	err := c.ShouldBind(&login)
	if err != nil {
		zap.L().Error("User Login 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 查询用户
	user, err := dao.GetUserByEmail(login.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrUserNotFound, c)
		} else {
			zap.L().Error("User Login 查询失败:" + err.Error())
			response.ErrWithMsg(code.InternalError, c)
		}
		return
	}
	// 检查密码
	if !utils.CheckPassword(user.Password, login.Password) {
		response.ErrWithMsg(code.ErrPassword, c)
		return
	}
	// 检查是否被封禁
	if user.Status == 2 {
		response.ErrWithMsg(code.Forbidden, c)
		return
	}
	// 生成token
	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		zap.L().Error("User Login 生成Token失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{
		"token":    token,
		"email":    user.Email,
		"nickname": user.Nickname,
		"id":       user.ID,
	}, c)
}

// 获取所有分类
func GetCategories(c *gin.Context) {
	categories, err := dao.GetAllCategories()
	if err != nil {
		zap.L().Error("GetCategories:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(categories, c)
}

// 根据分类获取文章
func GetCategoryArticles(c *gin.Context) {
	var q dto.CategoryArticlesQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		zap.L().Error("GetCategoryArticles 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if q.Page < 1 {
		q.Page = 1
	}
	articles, err := dao.GetArticleByCategory(q.CategoryID, q.Page, 10)
	if err != nil {
		zap.L().Error("GetCategoryArticles:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articles, c)
}

// 获取所有已使用的标签
func GetTags(c *gin.Context) {
	tags, err := dao.GetAllTags()
	if err != nil {
		zap.L().Error("GetTags:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(tags, c)
}

// 点赞文章（直接+1，不限制次数）
func LikeArticle(c *gin.Context) {
	var req dto.ArticleLikeReq
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("LikeArticle 参数错误:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	err = dao.IncrementLikeCount(req.ArticleID)
	if err != nil {
		zap.L().Error("LikeArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("点赞成功", c)
}
