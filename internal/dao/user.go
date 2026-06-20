package dao

import (
	"blog/core"
	"blog/models"

	"go.uber.org/zap"
)

// 根据邮箱查询用户是否存在
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := core.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// 创建用户
func CreateUser(user models.User) error {
	err := core.DB.Create(&user).Error
	if err != nil {
		zap.L().Error("CreateUser()" + err.Error())
		return err
	}
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint64) (models.User, error) {
	var user models.User
	err := core.DB.First(&user, id).Error
	if err != nil {
		zap.L().Error("GetUserByID:" + err.Error())
		return user, err
	}
	return user, nil
}

// GetUsersByPage 分页获取用户列表
func GetUsersByPage(page int, pageSize int) ([]models.User, int64, error) {
	var users []models.User = make([]models.User, 0)
	var total int64
	err := core.DB.Model(&models.User{}).Count(&total).Error
	if err != nil {
		zap.L().Error("GetUsersByPage count:" + err.Error())
		return users, total, err
	}
	err = core.DB.Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&users).Error
	if err != nil {
		zap.L().Error("GetUsersByPage:" + err.Error())
		return users, total, err
	}
	return users, total, nil
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(id uint64, status uint64) error {
	err := core.DB.Model(&models.User{}).Where("id = ?", id).
		Update("status", status).Error
	if err != nil {
		zap.L().Error("UpdateUserStatus:" + err.Error())
		return err
	}
	return nil
}

// DeleteUserByID 删除用户
func DeleteUserByID(id uint64) error {
	err := core.DB.Delete(&models.User{}, id).Error
	if err != nil {
		zap.L().Error("DeleteUserByID:" + err.Error())
		return err
	}
	return nil
}
