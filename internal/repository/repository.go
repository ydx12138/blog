package repository

import (
	"blog/models"
	"blog/models/vo"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// 根据email修改用户密码
func (r *Repository) UpdateUserPassword(email, password string) error {
	err := r.db.Model(&models.User{}).
		Where("email = ?", email).
		Update("password", password).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetArticleByPage(page int, pageSize int) ([]vo.ArticleSimple, int64, error) {
	articleList := make([]vo.ArticleSimple, 0)
	var total int64
	err := r.db.Model(&models.Article{}).Where("status = ?", 2).Count(&total).Error
	if err != nil {
		zap.L().Error("GetArticleByPage count:" + err.Error())
		return articleList, total, err
	}
	err = r.db.Model(models.Article{}).
		Select(`
			article.id, article.title, article.summary, article.cover,
			article.category_id, c.name AS category_name,
			article.view_count, article.like_count, article.comment_count,
			article.tags, article.created_at, article.updated_at
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Order("article.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articleList).Error
	return articleList, total, err
}

func (r *Repository) GetArticleDetail(id uint64) (vo.ArticleDetail, error) {
	var detail vo.ArticleDetail
	result := r.db.Model(models.Article{}).
		Select(`
			article.id, article.title, article.summary, article.content,
			article.cover, c.name AS category_name, article.view_count,
			article.like_count, article.comment_count, article.publish_time,
			article.tags, article.content_type
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Where("article.id = ?", id).
		Scan(&detail)
	if result.Error != nil {
		return detail, result.Error
	}
	if result.RowsAffected == 0 {
		return detail, gorm.ErrRecordNotFound
	}
	return detail, nil
}

func (r *Repository) SearchArticleByKey(keyword string) ([]vo.ArticleSimple, error) {
	articleList := make([]vo.ArticleSimple, 0)
	err := r.db.Model(models.Article{}).
		Select(`
			article.id, article.title, article.summary, article.cover,
			article.category_id, c.name AS category_name,
			article.view_count, article.like_count, article.comment_count,
			article.tags, article.created_at, article.updated_at
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Where("article.title like ? or article.summary like ? or article.content like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&articleList).Error
	return articleList, err
}

func (r *Repository) GetArticleByCategory(categoryID uint64, page int, pageSize int) ([]vo.ArticleSimple, error) {
	articleList := make([]vo.ArticleSimple, 0)
	err := r.db.Model(models.Article{}).
		Select(`
			article.id, article.title, article.summary, article.cover,
			article.category_id, c.name AS category_name,
			article.view_count, article.like_count, article.comment_count,
			article.tags, article.created_at, article.updated_at
		`).
		Joins("LEFT JOIN category c ON article.category_id = c.id").
		Where("article.status = ?", 2).
		Where("article.category_id = ?", categoryID).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articleList).Error
	return articleList, err
}

func (r *Repository) IncrementViewCount(id uint64) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", r.db.Raw("view_count + ?", 1)).Error
}

func (r *Repository) AdminGetArticles(page int, pageSize int, keyword string, status int8) ([]models.Article, int64, error) {
	articles := make([]models.Article, 0)
	var total int64
	query := r.db.Model(&models.Article{}).Preload("Category")
	if keyword != "" {
		query = query.Where("title like ? or summary like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return articles, total, err
	}
	err := query.Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&articles).Error
	return articles, total, err
}

func (r *Repository) GetArticleByID(id uint64) (models.Article, error) {
	var article models.Article
	err := r.db.Preload("Category").First(&article, id).Error
	return article, err
}

func (r *Repository) CreateArticle(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *Repository) UpdateArticle(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *Repository) DeleteArticle(id uint64) error {
	return r.db.Delete(&models.Article{}, id).Error
}

func (r *Repository) GetDrafts(page int, pageSize int) ([]models.Article, int64, error) {
	articles := make([]models.Article, 0)
	var total int64
	if err := r.db.Model(&models.Article{}).Preload("Category").Where("status = ?", 1).Count(&total).Error; err != nil {
		return articles, total, err
	}
	err := r.db.Preload("Category").Where("status = ?", 1).
		Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&articles).Error
	return articles, total, err
}

func (r *Repository) GetAllTags() ([]string, error) {
	var articles []models.Article
	if err := r.db.Select("tags").Where("status = ? AND tags != ''", 2).Find(&articles).Error; err != nil {
		return nil, err
	}
	tagSet := make(map[string]struct{})
	for _, a := range articles {
		for _, t := range strings.Split(a.Tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tagSet[t] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(tagSet))
	for t := range tagSet {
		result = append(result, t)
	}
	return result, nil
}

func (r *Repository) IncrementLikeCount(articleID uint64) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", r.db.Raw("like_count + ?", 1)).Error
}

func (r *Repository) UpdateArticleCommentCount(articleID uint64, delta int) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("comment_count", r.db.Raw("comment_count + ?", delta)).Error
}

func (r *Repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *Repository) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *Repository) GetUsersByPage(page int, pageSize int, keyword string, status uint64) ([]models.User, int64, error) {
	users := make([]models.User, 0)
	var total int64
	query := r.db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("email like ? OR nickname like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return users, total, err
	}
	err := query.Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error
	return users, total, err
}

func (r *Repository) UpdateUserStatus(id uint64, status uint64) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) DeleteUserByID(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *Repository) LoginVerification(username, password string) (models.Admin, error) {
	var ad models.Admin
	if err := r.db.Where("username = ?", username).First(&ad).Error; err != nil {
		return ad, err
	}
	//临时应急
	if ad.Password != password {
		return models.Admin{}, errors.New("password error")
	}
	//if !utils.CheckPassword(ad.Password, password) {
	//	return models.Admin{}, errors.New("password error")
	//}
	return ad, nil
}

func (r *Repository) GetAllCategories() ([]models.Category, error) {
	categories := make([]models.Category, 0)
	err := r.db.Order("sort DESC").Find(&categories).Error
	return categories, err
}

func (r *Repository) GetOrCreateDefaultCategory() (models.Category, error) {
	var cat models.Category
	if err := r.db.Where("name = ?", "杂谈").First(&cat).Error; err == nil {
		return cat, nil
	}
	cat = models.Category{Name: "杂谈", Description: "未分类的杂谈文章", Sort: 0}
	err := r.db.Create(&cat).Error
	return cat, err
}

func (r *Repository) CreateComment(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *Repository) GetCommentsByArticle(articleID uint64, page int, pageSize int) ([]vo.CommentVO, int64, error) {
	comments := make([]vo.CommentVO, 0)
	var total int64
	if err := r.db.Model(&models.Comment{}).Where("article_id = ? AND status = ?", articleID, 1).Count(&total).Error; err != nil {
		return comments, total, err
	}
	err := r.db.Model(&models.Comment{}).
		Select(`
			comment.id, comment.article_id, a.title AS article_title,
			comment.user_id, u.nickname, comment.content,
			comment.parent_id, comment.status, comment.created_at
		`).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id").
		Where("comment.article_id = ? AND comment.status = ?", articleID, 1).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	return comments, total, err
}

func (r *Repository) GetAllComments(page int, pageSize int, keyword string, searchType string) ([]vo.CommentVO, int64, error) {
	comments := make([]vo.CommentVO, 0)
	var total int64
	query := r.db.Model(&models.Comment{}).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id")
	if keyword != "" {
		if searchType == "nickname" {
			query = query.Where("u.nickname like ?", "%"+keyword+"%")
		} else {
			query = query.Where("comment.content like ?", "%"+keyword+"%")
		}
	}
	if err := query.Count(&total).Error; err != nil {
		return comments, total, err
	}
	err := query.Select(`
			comment.id, comment.article_id, a.title AS article_title,
			comment.user_id, u.nickname, comment.content,
			comment.parent_id, comment.status, comment.created_at
		`).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	return comments, total, err
}

func (r *Repository) GetPendingComments(page int, pageSize int) ([]vo.CommentVO, int64, error) {
	comments := make([]vo.CommentVO, 0)
	var total int64
	if err := r.db.Model(&models.Comment{}).Where("status = ?", 3).Count(&total).Error; err != nil {
		return comments, total, err
	}
	err := r.db.Model(&models.Comment{}).
		Select(`
			comment.id, comment.article_id, a.title AS article_title,
			comment.user_id, u.nickname, comment.content,
			comment.parent_id, comment.status, comment.created_at
		`).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id").
		Where("comment.status = ?", 3).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	return comments, total, err
}

func (r *Repository) UpdateCommentStatus(id uint64, status int8) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) GetCommentByID(id uint64) (models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	return comment, err
}

func (r *Repository) DeleteComment(id uint64) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

type DashboardData struct {
	TotalArticles     int64 `json:"total_articles"`
	PublishedArticles int64 `json:"published_articles"`
	DraftArticles     int64 `json:"draft_articles"`
	TotalComments     int64 `json:"total_comments"`
	PendingComments   int64 `json:"pending_comments"`
	TotalUsers        int64 `json:"total_users"`
	TotalViews        int64 `json:"total_views"`
}

func (r *Repository) GetDashboard() (DashboardData, error) {
	var data DashboardData
	if err := r.db.Model(&models.Article{}).Count(&data.TotalArticles).Error; err != nil {
		return data, err
	}
	if err := r.db.Model(&models.Article{}).Where("status = ?", 2).Count(&data.PublishedArticles).Error; err != nil {
		return data, err
	}
	if err := r.db.Model(&models.Article{}).Where("status = ?", 1).Count(&data.DraftArticles).Error; err != nil {
		return data, err
	}
	if err := r.db.Model(&models.Comment{}).Count(&data.TotalComments).Error; err != nil {
		return data, err
	}
	if err := r.db.Model(&models.Comment{}).Where("status = ?", 3).Count(&data.PendingComments).Error; err != nil {
		return data, err
	}
	if err := r.db.Model(&models.User{}).Count(&data.TotalUsers).Error; err != nil {
		return data, err
	}
	err := r.db.Model(&models.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&data.TotalViews).Error
	return data, err
}

func (r *Repository) SaveToken(userID uint64, token string) error {
	tk := models.Token{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return r.db.Where("user_id = ?", userID).Assign(tk).FirstOrCreate(&tk).Error
}
