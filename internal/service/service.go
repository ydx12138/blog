package service

import (
	"blog/internal/repository"
	"blog/internal/utils"
	"blog/models"
	"blog/models/dto"
	"blog/models/vo"
	"context"
	"crypto/subtle"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 验证码过期时间
const registerCodeTTL = 60 * time.Second
const resetPasswordCodeTTL = 60 * time.Second

var (
	ErrUserExists       = errors.New("user exists")
	ErrUserNotExists    = errors.New("user not exists")
	ErrPassword         = errors.New("password error")
	ErrUserDisabled     = errors.New("user disabled")
	ErrVerificationCode = errors.New("verification code invalid or expired")
	ErrCoderepeated     = errors.New("verification codes cannot be obtained repeatedly")
)

type Service struct {
	repo  *repository.Repository
	redis *redis.Client
}

func New(repo *repository.Repository, redis *redis.Client) *Service {
	return &Service{repo: repo, redis: redis}
}

// refreshToken是否还在
func (s *Service) RefreshTokenIsExist(userid string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uid, err := strconv.Atoi(userid)
	if err != nil {
		zap.L().Error("Invalid user id" + err.Error())
		return false
	}
	//
	refreshToken := s.redis.Get(ctx, refreshTokenKey(uint64(uid))).Val()
	if refreshToken == "" {
		return false
	}
	return true

}

// 修改电话号码
func (s *Service) UpdatePhoneNumber(email, phone, verificationcode string) error {
	//这个手机号是否已经被其他账号绑定
	//
	return nil
}

func (s *Service) UpdatePasswordByCode(email, password string, code string) error {
	//email对应的验证码是否存在
	result, err := s.CheckCode(email)
	if err != nil {
		return err
	}
	if result == false {
		return ErrVerificationCode
	}
	//检查重置密码验证码
	res, err := s.VerifyResetPasswordCode(email, code)
	if err != nil {
		return err
	}
	if res == false {
		return ErrVerificationCode
	}
	//修改密码
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	if err = s.repo.UpdateUserPassword(email, hashPassword); err != nil {
		return err
	}
	return nil
}

// 验证重置密码的验证码是否正确
func (s *Service) VerifyResetPasswordCode(email string, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//检查验证码是否存在
	result, err := s.redis.ZRangeWithScores(ctx, resetPasswordCodeKey(email), 0, -1).Result()
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, ErrVerificationCode
	}
	//是否一样,如果不一样，次数减一
	for _, value := range result {
		if value.Member != code {
			err = s.redis.ZIncrBy(ctx, resetPasswordCodeKey(email), -1, value.Member.(string)).Err()
			if err != nil {
				return false, err
			}
			if err = s.DeleteCodeEfftive(email); err != nil {
				return false, err
			}
			return false, ErrVerificationCode
		}
	}
	//如果次数-1之后为0，删除这个key

	//验证码通过
	return true, nil
}

// 如果次数为0，删除验证码
func (s *Service) DeleteCodeEfftive(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := s.redis.ZRangeWithScores(ctx, resetPasswordCodeKey(email), 0, -1).Result()
	if err != nil {
		return err
	}
	for _, value := range res {
		if value.Score == 0 {
			err = s.redis.Del(ctx, resetPasswordCodeKey(email)).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 给忘记密码的邮箱发验证码
func (s *Service) SendCodeForgetPwd(email string) error {
	email = normalizeEmail(email)
	//邮箱是否已经存在，如果不存在，返回
	user, err := s.repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if user.ID == 0 {
		return ErrUserNotExists
	}
	//上次已发的验证码是否还在有效期内
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := s.CheckCode(email)
	if err != nil {
		return err
	}
	if result == true {
		return ErrCoderepeated
	}
	//生成验证码
	verifyCode, err := utils.GenerateCode()
	if err != nil {
		return err
	}
	//存到redis
	if err = s.SaveCodeForgetPwd(email, verifyCode, 5); err != nil {
		return err
	}
	//发验证码，如果发送失败，把刚刚存的验证码删除
	key := resetPasswordCodeKey(email)
	if err = utils.SendEmailToQQ(email, "YDX Blog 重置密码验证码", verifyCode); err != nil {
		_ = s.redis.Del(ctx, key).Err()
		return err
	}
	return nil
}

// 保存重置密码验证码
func (s *Service) SaveCodeForgetPwd(email, verifyCode string, effective float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := resetPasswordCodeKey(email)
	//存到zset,5次访问机会
	if err := s.redis.ZAdd(ctx, key, redis.Z{Score: effective, Member: verifyCode}).Err(); err != nil {
		return err
	}
	//60秒过期
	err := s.redis.Expire(ctx, key, resetPasswordCodeTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

// 检查是否已经存在重置密码验证码
func (s *Service) CheckCode(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := resetPasswordCodeKey(email)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Service) GetArticles(page, pageSize int) ([]vo.ArticleSimple, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.GetArticleByPage(page, pageSize)
}

func (s *Service) GetArticle(id uint64) (vo.ArticleDetail, error) {
	detail, err := s.repo.GetArticleDetail(id)
	if err != nil {
		return detail, err
	}
	_ = s.repo.IncrementViewCount(id)
	return detail, nil
}

func (s *Service) SearchArticle(keyword string) ([]vo.ArticleSimple, error) {
	return s.repo.SearchArticleByKey(keyword)
}

func (s *Service) SendRegisterCode(req dto.SendRegisterCodeReq) error {
	email := normalizeEmail(req.Email)
	//查找邮箱
	existUser, err := s.repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	//如果邮箱存在
	if existUser.ID != 0 {
		return ErrUserExists
	}
	//如果redis客户端为nil
	if s.redis == nil {
		return errors.New("redis client is nil")
	}
	//生成验证码
	verifyCode, err := utils.GenerateCode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//key:verifyCode
	key := registerCodeKey(email)
	//验证码存redis
	if err := s.redis.Set(ctx, key, verifyCode, registerCodeTTL).Err(); err != nil {
		return err
	}
	//发验证码
	if err := utils.SendEmailToQQ(email, "YDX Blog 注册验证码", verifyCode); err != nil {
		_ = s.redis.Del(ctx, key).Err()
		return err
	}
	return nil
}

// 注册
func (s *Service) Register(req dto.UserRegister) error {
	email := normalizeEmail(req.Email)
	//判断验证码是否无误
	if err := s.verifyRegisterCode(email, req.Code); err != nil {
		return err
	}
	//用户是否已经存在
	existUser, err := s.repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existUser.ID != 0 {
		return ErrUserExists
	}
	//加密密码
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	//创建新用户
	if err := s.repo.CreateUser(models.User{
		Email:    email,
		Password: hashPassword,
		Nickname: req.Nickname,
	}); err != nil {
		return err
	}
	//删除验证码
	_ = s.deleteRegisterCode(email)
	return nil
}

// 登录
func (s *Service) UserLogin(req dto.UserLogin) (map[string]interface{}, error) {
	//根据邮箱查用户
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	//密码是否正确
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, ErrPassword
	}
	//用户是否被禁用
	if user.Status == 2 {
		return nil, ErrUserDisabled
	}
	//生成refreshToken，放入用户ID和用户身份(user)
	refreshToken, err := utils.GenerateUserToken(user.ID, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, err
	}
	//生成accessToken，放入用户ID和用户身份(user)
	accessToken, err := utils.GenerateUserToken(user.ID, 15*time.Minute, "access")
	if err != nil {
		return nil, err
	}
	//保存refreshToken
	if err = s.SaveRefreshToken(refreshTokenKey(user.ID), refreshToken, 7*24*time.Hour); err != nil {
		return nil, err
	}
	//accessToken和refreshToken全部返回
	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"email":         user.Email,
		"nickname":      user.Nickname,
		"id":            user.ID,
	}, nil
}

// 把refreshToken存入redis
func (s *Service) SaveRefreshToken(key, token string, duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.redis.Set(ctx, key, token, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *Service) GetCategoryArticles(categoryID uint64, page int) ([]vo.ArticleSimple, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.GetArticleByCategory(categoryID, page, 10)
}

func (s *Service) GetTags() ([]string, error) {
	return s.repo.GetAllTags()
}

func (s *Service) LikeArticle(articleID uint64) error {
	return s.repo.IncrementLikeCount(articleID)
}

func (s *Service) GetComments(articleID uint64, page int) ([]vo.CommentVO, int64, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.GetCommentsByArticle(articleID, page, 10)
}

// 保存评论
func (s *Service) CreateComment(req dto.CreateCommentReq, userID uint64) error {
	//处理敏感词
	var words []string
	if utils.Has(req.Content) == true {
		words = utils.FindAll(req.Content)
		req.Content = utils.Replace(req.Content)
	}

	comment := models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Status:    1,
		HitWords:  strings.Join(words, ","),
	}

	if err := s.repo.CreateComment(&comment); err != nil {
		return err
	}
	_ = s.repo.UpdateArticleCommentCount(comment.ArticleID, 1)
	return nil
}

func (s *Service) AdminLogin(req dto.AdminLogin) (map[string]interface{}, error) {
	admin, err := s.repo.LoginVerification(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateAdminToken(admin.ID, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"token":    token,
		"nickname": admin.Nickname,
		"username": admin.Username,
	}, nil
}

func (s *Service) Dashboard() (repository.DashboardData, error) {
	return s.repo.GetDashboard()
}

func (s *Service) AdminArticles(q dto.AdminArticleQuery) ([]models.Article, int64, error) {
	page, pageSize := normalizePage(q.Page, q.PageSize)
	return s.repo.AdminGetArticles(page, pageSize, q.Keyword, q.Status)
}

func (s *Service) AdminArticle(id uint64) (models.Article, error) {
	return s.repo.GetArticleByID(id)
}

func (s *Service) CreateArticle(req dto.CreateArticleReq, authorID uint64) error {
	categoryID := req.CategoryID
	if categoryID == 0 {
		defaultCat, err := s.repo.GetOrCreateDefaultCategory()
		if err == nil {
			categoryID = defaultCat.ID
		}
	}
	article := models.Article{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		ContentType: req.ContentType,
		Cover:       req.Cover,
		CategoryID:  categoryID,
		Tags:        req.Tags,
		Status:      req.Status,
		AuthorID:    authorID,
	}
	if req.Status == 2 {
		now := time.Now()
		article.PublishTime = &now
	}
	return s.repo.CreateArticle(&article)
}

func (s *Service) UpdateArticle(req dto.UpdateArticleReq) error {
	article, err := s.repo.GetArticleByID(req.ID)
	if err != nil {
		return err
	}
	article.Title = req.Title
	article.Summary = req.Summary
	article.Content = req.Content
	article.ContentType = req.ContentType
	article.Cover = req.Cover
	article.CategoryID = req.CategoryID
	article.Tags = req.Tags
	article.Status = req.Status
	if req.Status == 2 && article.PublishTime == nil {
		now := time.Now()
		article.PublishTime = &now
	}
	return s.repo.UpdateArticle(&article)
}

func (s *Service) DeleteArticle(id uint64) error {
	return s.repo.DeleteArticle(id)
}

func (s *Service) Drafts(q dto.AdminArticleQuery) ([]models.Article, int64, error) {
	page, pageSize := normalizePage(q.Page, q.PageSize)
	return s.repo.GetDrafts(page, pageSize)
}

func (s *Service) PublishArticle(id uint64) error {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		return err
	}
	article.Status = 2
	now := time.Now()
	article.PublishTime = &now
	return s.repo.UpdateArticle(&article)
}

func (s *Service) AllComments(page, pageSize int, keyword string, searchType string) ([]vo.CommentVO, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	return s.repo.GetAllComments(page, pageSize, keyword, searchType)
}

func (s *Service) PendingComments(page, pageSize int) ([]vo.CommentVO, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	return s.repo.GetPendingComments(page, pageSize)
}

func (s *Service) SetCommentStatus(id uint64, status int8) error {
	old, err := s.repo.GetCommentByID(id)
	if err != nil {
		return err
	}
	if old.Status == 1 && status != 1 {
		_ = s.repo.UpdateArticleCommentCount(old.ArticleID, -1)
	}
	if old.Status != 1 && status == 1 {
		_ = s.repo.UpdateArticleCommentCount(old.ArticleID, 1)
	}
	return s.repo.UpdateCommentStatus(id, status)
}

func (s *Service) DeleteComment(id uint64) error {
	comment, err := s.repo.GetCommentByID(id)
	if err == nil && comment.Status == 1 {
		_ = s.repo.UpdateArticleCommentCount(comment.ArticleID, -1)
	}
	return s.repo.DeleteComment(id)
}

func (s *Service) Users(page, pageSize int, keyword string, status uint64) ([]models.User, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	return s.repo.GetUsersByPage(page, pageSize, keyword, status)
}

func (s *Service) BanUser(id uint64) error {
	return s.repo.UpdateUserStatus(id, 2)
}

func (s *Service) UnbanUser(id uint64) error {
	return s.repo.UpdateUserStatus(id, 1)
}

func (s *Service) DeleteUser(id uint64) error {
	return s.repo.DeleteUserByID(id)
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// 邮箱对应的验证码是否正确?
func (s *Service) verifyRegisterCode(email string, input string) error {
	if s.redis == nil {
		return errors.New("redis client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stored, err := s.redis.Get(ctx, registerCodeKey(email)).Result()
	if errors.Is(err, redis.Nil) {
		return ErrVerificationCode
	}
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	if subtle.ConstantTimeCompare([]byte(stored), []byte(input)) != 1 {
		return ErrVerificationCode
	}
	return nil
}

// 删除邮箱对应的验证码
func (s *Service) deleteRegisterCode(email string) error {
	if s.redis == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.redis.Del(ctx, registerCodeKey(email)).Err()
}

func registerCodeKey(email string) string {
	return "register:email_code:" + normalizeEmail(email)
}
func resetPasswordCodeKey(email string) string {
	return "reset_password:email_code:" + normalizeEmail(email)
}

// refreshToken存入redis时的key
func refreshTokenKey(userid uint64) string {
	return strconv.Itoa(int(userid)) + ":refreshToken:"
}

// 把邮件规范化（大写改成小写）
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
