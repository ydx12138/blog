package dao

import (
	"blog/core"
	"blog/models"
	"time"
)

// SaveToken 保存用户token
func SaveToken(userID uint64, token string) error {
	tk := models.Token{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return core.DB.Where("user_id = ?", userID).Assign(tk).FirstOrCreate(&tk).Error
}

// DeleteToken 删除用户token（下线）
func DeleteToken(userID uint64) error {
	return core.DB.Where("user_id = ?", userID).Delete(&models.Token{}).Error
}

// TokenExists 检查token是否有效
func TokenExists(userID uint64) bool {
	var count int64
	core.DB.Model(&models.Token{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&count)
	return count > 0
}
