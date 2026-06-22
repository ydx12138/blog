package dao

import (
	"blog/core"
	"blog/models"
	"blog/models/vo"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(comment *models.Comment) error {
	err := core.DB.Create(comment).Error
	if err != nil {
		zap.L().Error("CreateComment:" + err.Error())
		return err
	}
	return nil
}

// GetCommentsByArticle 获取文章评论列表（已审核通过的）
func GetCommentsByArticle(articleID uint64, page int, pageSize int) ([]vo.CommentVO, int64, error) {
	var comments []vo.CommentVO = make([]vo.CommentVO, 0)
	var total int64
	err := core.DB.Model(&models.Comment{}).
		Where("article_id = ? AND status = ?", articleID, 1).Count(&total).Error
	if err != nil {
		zap.L().Error("GetCommentsByArticle count:" + err.Error())
		return comments, total, err
	}
	err = core.DB.
		Model(&models.Comment{}).
		Select(`
			comment.id,
			comment.article_id,
			a.title AS article_title,
			comment.user_id,
			u.nickname,
			comment.content,
			comment.parent_id,
			comment.status,
			comment.created_at
		`).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id").
		Where("comment.article_id = ? AND comment.status = ?", articleID, 1).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	if err != nil {
		zap.L().Error("GetCommentsByArticle:" + err.Error())
		return comments, total, err
	}
	return comments, total, nil
}

// GetAllComments 获取全部评论（按时间倒序，含分页，支持搜索）
func GetAllComments(page int, pageSize int, keyword string, searchType string) ([]vo.CommentVO, int64, error) {
	var comments []vo.CommentVO = make([]vo.CommentVO, 0)
	var total int64
	query := core.DB.Model(&models.Comment{}).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id")
	if keyword != "" {
		if searchType == "nickname" {
			query = query.Where("u.nickname like ?", "%"+keyword+"%")
		} else {
			query = query.Where("comment.content like ?", "%"+keyword+"%")
		}
	}
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("GetAllComments count:" + err.Error())
		return comments, total, err
	}
	err = query.
		Select(`
			comment.id,
			comment.article_id,
			a.title AS article_title,
			comment.user_id,
			u.nickname,
			comment.content,
			comment.parent_id,
			comment.status,
			comment.created_at
		`).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	if err != nil {
		zap.L().Error("GetAllComments:" + err.Error())
		return comments, total, err
	}
	return comments, total, nil
}

// GetPendingComments 获取待审核评论
func GetPendingComments(page int, pageSize int) ([]vo.CommentVO, int64, error) {
	var comments []vo.CommentVO = make([]vo.CommentVO, 0)
	var total int64
	err := core.DB.Model(&models.Comment{}).
		Where("status = ?", 3).Count(&total).Error
	if err != nil {
		zap.L().Error("GetPendingComments count:" + err.Error())
		return comments, total, err
	}
	err = core.DB.
		Model(&models.Comment{}).
		Select(`
			comment.id,
			comment.article_id,
			a.title AS article_title,
			comment.user_id,
			u.nickname,
			comment.content,
			comment.parent_id,
			comment.status,
			comment.created_at
		`).
		Joins("LEFT JOIN user u ON comment.user_id = u.id").
		Joins("LEFT JOIN article a ON comment.article_id = a.id").
		Where("comment.status = ?", 3).
		Order("comment.created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&comments).Error
	if err != nil {
		zap.L().Error("GetPendingComments:" + err.Error())
		return comments, total, err
	}
	return comments, total, nil
}

// UpdateCommentStatus 更新评论状态
func UpdateCommentStatus(id uint64, status int8) error {
	err := core.DB.Model(&models.Comment{}).Where("id = ?", id).
		Update("status", status).Error
	if err != nil {
		zap.L().Error("UpdateCommentStatus:" + err.Error())
		return err
	}
	return nil
}

// GetCommentByID 获取单条评论
func GetCommentByID(id uint64) (models.Comment, error) {
	var c models.Comment
	err := core.DB.First(&c, id).Error
	if err != nil {
		zap.L().Error("GetCommentByID:" + err.Error())
		return c, err
	}
	return c, nil
}

// DeleteComment 删除评论
func DeleteComment(id uint64) error {
	err := core.DB.Delete(&models.Comment{}, id).Error
	if err != nil {
		zap.L().Error("DeleteComment:" + err.Error())
		return err
	}
	return nil
}
