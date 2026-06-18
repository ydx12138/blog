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
