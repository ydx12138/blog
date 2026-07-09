package handler

import (
	"blog/config"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/models/dto"
	"blog/pkg/code"
	"blog/pkg/response"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	User  *UserHandler
	Admin *AdminHandler
}

// 其它地方访问不到svc未导出字段，只有结构体方法能够访问到svc，从而访问到service层方法
type UserHandler struct {
	svc *service.Service
}

type AdminHandler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		User:  &UserHandler{svc: svc},
		Admin: &AdminHandler{svc: svc},
	}
}

// 刷新token
func (h *UserHandler) TokenRefresh(c *gin.Context) {
	//从context中得到token
	token := utils.GetTokenFromContext(c)
	if token == "" {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//解析token，得到Data
	data, err := utils.GetDataFromToken(token)
	if data == nil || err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//解析token，得到claim
	var claim *utils.CustomClaims
	if claim = utils.GetClaimFromData(data); claim == nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//如果type==refresh,有效，且redis里存在，则创建新accessToken，Abort+return
	if claim.Type == "refresh" && data.Valid && h.svc.RefreshTokenIsExist(strconv.FormatUint(claim.UserID, 10)) == true {
		accessToken, err := utils.GenerateUserToken(claim.UserID, 15*time.Minute, "access")
		if err != nil {
			zap.L().Error("generate access token failed" + err.Error())
			response.ErrWithMsg(code.InternalError, c)
			return
		}

		response.SuccessWithData(map[string]interface{}{"access_token": accessToken}, c)
		return
	}
	//如果type==refresh,无效，或者redis里不存在，则401，Abort+return
	if claim.Type == "refresh" && (data.Valid == false && h.svc.RefreshTokenIsExist(claim.ID) == false) {
		response.ErrWithMsg(code.RefreshTokenExpired, c)
		return
	}
}

// 主要：验证token过期否
func (h *UserHandler) UsersMe(c *gin.Context) {
	return
}

// 修改手机号
func (h *UserHandler) UpdatePhoneNumber(c *gin.Context) {
	var q dto.UserPutPhone
	if err := c.ShouldBind(&q); err != nil {
		zap.L().Error("ForgetPassword" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.UpdatePhoneNumber(q.Email, q.Phone, q.VerificationCode); err != nil {

	}
	response.ErrWithMsg(code.Success, c)
}

// 重置密码--1.发验证码
func (h *UserHandler) SendCodeForgetPwd(c *gin.Context) {
	var q dto.SendRegisterCodeReq
	err := c.ShouldBind(&q)
	if err != nil {
		zap.L().Error("ForgetPassword" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//发验证码
	if err = h.svc.SendCodeForgetPwd(q.Email); err != nil {
		zap.L().Error("ForgetPassword:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("验证码发送成功", c)
}

// 重置密码--2.修改密码
func (h *UserHandler) UpdatePasswordByCode(c *gin.Context) {
	//参数
	var q dto.UserUpdatePassword
	err := c.ShouldBind(&q)
	if err != nil {
		zap.L().Error("UpdatePasswordByCode" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//修改密码
	if err = h.svc.UpdatePasswordByCode(q.Email, q.Password, q.Code); err != nil {
		zap.L().Error("UpdatePasswordByCode:" + err.Error())
		response.ErrWithMsg(code.ErrorMsg(err), c)
		return
	}
	response.SuccessWithMsg("密码重置成功", c)
}

func (h *UserHandler) GetArticles(c *gin.Context) {
	var q dto.PageQueryWithSize
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page = 1
	}
	articles, total, err := h.svc.GetArticles(q.Page, q.PageSize)
	if err != nil {
		zap.L().Error("GetArticles:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": articles, "total": total}, c)
}

func (h *UserHandler) GetArticle(c *gin.Context) {
	var q struct {
		ID uint64 `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	article, err := h.svc.GetArticle(q.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrArticleNotFound, c)
			return
		}
		zap.L().Error("GetArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(article, c)
}

func (h *UserHandler) SearchArticle(c *gin.Context) {
	var q dto.ArticleKeyWord
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	articles, err := h.svc.SearchArticle(q.Keyword)
	if err != nil {
		zap.L().Error("SearchArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articles, c)
}

// 注册
func (h *UserHandler) Register(c *gin.Context) {
	//参数
	var req dto.UserRegister
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("Register:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//注册
	if err := h.svc.Register(req); err != nil {
		//用户已存在
		if errors.Is(err, service.ErrUserExists) {
			response.ErrWithMsg(code.ErrUserExist, c)
			return
		}
		//验证码无效
		if errors.Is(err, service.ErrVerificationCode) {
			response.ErrWithMsg(code.ErrVerificationCode, c)
			return
		}
		//其它原因
		zap.L().Error("Register:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("注册成功", c)
}

// 发送验证码
func (h *UserHandler) SendRegisterCode(c *gin.Context) {
	var req dto.SendRegisterCodeReq
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("SendRegisterCode:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//
	if err := h.svc.SendRegisterCode(req); err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.ErrWithMsg(code.ErrUserExist, c)
			return
		}
		zap.L().Error("SendRegisterCode:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("验证码已发送", c)
}

// 登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLogin
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//登录
	data, err := h.svc.UserLogin(req)
	if err != nil {
		switch {
		//用户不存在
		case errors.Is(err, gorm.ErrRecordNotFound):
			response.ErrWithMsg(code.ErrUserNotFound, c)
		//密码错误
		case errors.Is(err, service.ErrPassword):
			response.ErrWithMsg(code.ErrPassword, c)
		//用户被禁用
		case errors.Is(err, service.ErrUserDisabled):
			response.ErrWithMsg(code.Forbidden, c)
		default:
			//服务器错误
			zap.L().Error("Login:" + err.Error())
			response.ErrWithMsg(code.InternalError, c)
		}
		return
	}
	response.SuccessWithData(data, c)
}

func (h *UserHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories()
	if err != nil {
		zap.L().Error("GetCategories:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(categories, c)
}

func (h *UserHandler) GetCategoryArticles(c *gin.Context) {
	var q dto.CategoryArticlesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	articles, err := h.svc.GetCategoryArticles(q.CategoryID, q.Page)
	if err != nil {
		zap.L().Error("GetCategoryArticles:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(articles, c)
}

func (h *UserHandler) GetTags(c *gin.Context) {
	tags, err := h.svc.GetTags()
	if err != nil {
		zap.L().Error("GetTags:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(tags, c)
}

func (h *UserHandler) LikeArticle(c *gin.Context) {
	var req dto.ArticleLikeReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.LikeArticle(req.ArticleID); err != nil {
		zap.L().Error("LikeArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("点赞成功", c)
}

func (h *UserHandler) GetComments(c *gin.Context) {
	var q dto.CommentListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	comments, total, err := h.svc.GetComments(q.ArticleID, q.Page)
	if err != nil {
		zap.L().Error("GetComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": comments, "total": total}, c)
}

func (h *UserHandler) CreateComment(c *gin.Context) {
	var req dto.CreateCommentReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	// 获取userid
	userID, ok := c.Get("userID")
	if !ok {
		response.ErrWithMsg(code.Unauthorized, c)
		return
	}
	uid, ok := userID.(uint64)
	if !ok {
		response.ErrWithMsg(code.Unauthorized, c)
		return
	}
	// 保存评论
	if err := h.svc.CreateComment(req, uid); err != nil {
		zap.L().Error("CreateComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("评论成功", c)
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req dto.AdminLogin
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	data, err := h.svc.AdminLogin(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrUserNotFound, c)
		} else if err.Error() == "password error" {
			response.ErrWithMsg(code.ErrPassword, c)
		} else {
			zap.L().Error("AdminLogin:" + err.Error())
			response.ErrWithMsg(code.InternalError, c)
		}
		return
	}
	response.SuccessWithData(data, c)
}

// 获取仪表盘
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	data, err := h.svc.Dashboard()
	if err != nil {
		zap.L().Error("GetDashboard:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(data, c)
}

// 获取文章列表
func (h *AdminHandler) GetArticles(c *gin.Context) {
	var q dto.AdminArticleQuery
	_ = c.ShouldBindQuery(&q)
	articles, total, err := h.svc.AdminArticles(q)
	if err != nil {
		zap.L().Error("AdminGetArticles:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": articles, "total": total}, c)
}

func (h *AdminHandler) GetArticle(c *gin.Context) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	article, err := h.svc.AdminArticle(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrArticleNotFound, c)
			return
		}
		zap.L().Error("AdminGetArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(article, c)
}

func (h *AdminHandler) CreateArticle(c *gin.Context) {
	var req dto.CreateArticleReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	authorID, ok := currentUserID(c)
	if !ok {
		response.ErrWithMsg(code.Unauthorized, c)
		return
	}
	if err := h.svc.CreateArticle(req, authorID); err != nil {
		zap.L().Error("CreateArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("创建成功", c)
}

func (h *AdminHandler) UpdateArticle(c *gin.Context) {
	var req dto.UpdateArticleReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if req.ID == 0 {
		id, err := parseParamID(c)
		if err != nil {
			response.ErrWithMsg(code.BadRequest, c)
			return
		}
		req.ID = id
	}
	if err := h.svc.UpdateArticle(req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrArticleNotFound, c)
			return
		}
		zap.L().Error("UpdateArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("更新成功", c)
}

func (h *AdminHandler) DeleteArticle(c *gin.Context) {
	id, err := idFromRequest(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.DeleteArticle(id); err != nil {
		zap.L().Error("DeleteArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}

func (h *AdminHandler) GetDrafts(c *gin.Context) {
	var q dto.AdminArticleQuery
	_ = c.ShouldBindQuery(&q)
	articles, total, err := h.svc.Drafts(q)
	if err != nil {
		zap.L().Error("GetDrafts:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": articles, "total": total}, c)
}

func (h *AdminHandler) PublishArticle(c *gin.Context) {
	id, err := idFromRequest(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.PublishArticle(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrArticleNotFound, c)
			return
		}
		zap.L().Error("PublishArticle:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("发布成功", c)
}

func (h *AdminHandler) GetAllComments(c *gin.Context) {
	var q dto.AdminCommentQuery
	_ = c.ShouldBindQuery(&q)
	comments, total, err := h.svc.AllComments(q.Page, q.PageSize, c.Query("keyword"), c.Query("type"))
	if err != nil {
		zap.L().Error("GetAllComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": comments, "total": total}, c)
}

func (h *AdminHandler) GetPendingComments(c *gin.Context) {
	var q dto.AdminCommentQuery
	_ = c.ShouldBindQuery(&q)
	comments, total, err := h.svc.PendingComments(q.Page, q.PageSize)
	if err != nil {
		zap.L().Error("GetPendingComments:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithData(map[string]interface{}{"list": comments, "total": total}, c)
}

func (h *AdminHandler) ApproveComment(c *gin.Context) {
	h.setCommentStatus(c, 1, "审核通过")
}

func (h *AdminHandler) RejectComment(c *gin.Context) {
	h.setCommentStatus(c, 3, "已驳回")
}

func (h *AdminHandler) SetCommentStatus(c *gin.Context) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if c.ShouldBind(&body) != nil || (body.Status != 1 && body.Status != 3) {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.SetCommentStatus(id, body.Status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrCommentNotFound, c)
			return
		}
		zap.L().Error("SetCommentStatus:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("状态已更新", c)
}

func (h *AdminHandler) DeleteComment(c *gin.Context) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.DeleteComment(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrCommentNotFound, c)
			return
		}
		zap.L().Error("DeleteComment:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	var q dto.AdminArticleQuery
	_ = c.ShouldBindQuery(&q)
	var status uint64
	if statusStr := c.Query("status"); statusStr != "" {
		status, _ = strconv.ParseUint(statusStr, 10, 64)
	}
	users, total, err := h.svc.Users(q.Page, q.PageSize, c.Query("keyword"), status)
	if err != nil {
		zap.L().Error("GetUsers:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	type safeUser struct {
		ID        uint64 `json:"id"`
		Email     string `json:"email"`
		Nickname  string `json:"nickname"`
		Status    uint64 `json:"status"`
		CreatedAt string `json:"created_at"`
	}
	safeUsers := make([]safeUser, 0, len(users))
	for _, u := range users {
		safeUsers = append(safeUsers, safeUser{
			ID:        u.ID,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.SuccessWithData(map[string]interface{}{"list": safeUsers, "total": total}, c)
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	h.setUserStatus(c, true)
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	h.setUserStatus(c, false)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.DeleteUser(id); err != nil {
		zap.L().Error("DeleteUser:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg("删除成功", c)
}

var allowedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

// 上传图片
func (h *AdminHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zap.L().Error("UploadImage:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			zap.L().Error("UploadImage:" + err.Error())
		}
	}(file)
	//图片大小10MB以下，且必须是指定的后缀
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] || header.Size > 10*1024*1024 {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	//图片名字
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(header.Filename, ext), ext)
	//上传，获得url
	url, err := utils.UploadToOss(file, config.Cfg.OssConfig.Image_path, filename)
	if err != nil {
		zap.L().Error("UploadImage:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	//返回url
	c.JSON(http.StatusOK, response.Response{Code: 0, Message: "上传成功", Data: map[string]string{"url": url}})
}

func (h *AdminHandler) setCommentStatus(c *gin.Context, status int8, msg string) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if err := h.svc.SetCommentStatus(id, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrWithMsg(code.ErrCommentNotFound, c)
			return
		}
		zap.L().Error("SetCommentStatus:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	response.SuccessWithMsg(msg, c)
}

func (h *AdminHandler) setUserStatus(c *gin.Context, banned bool) {
	id, err := parseParamID(c)
	if err != nil {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	if banned {
		err = h.svc.BanUser(id)
	} else {
		err = h.svc.UnbanUser(id)
	}
	if err != nil {
		zap.L().Error("SetUserStatus:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	if banned {
		response.SuccessWithMsg("封禁成功", c)
	} else {
		response.SuccessWithMsg("解封成功", c)
	}
}

func parseParamID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

func idFromRequest(c *gin.Context) (uint64, error) {
	if id, err := parseParamID(c); err == nil && id > 0 {
		return id, nil
	}
	var req dto.IDReq
	if err := c.ShouldBind(&req); err != nil {
		return 0, err
	}
	return req.ID, nil
}

func currentUserID(c *gin.Context) (uint64, bool) {
	userID, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	uid, ok := userID.(uint64)
	return uid, ok
}
