package dao

import (
	"blog/core"
	"blog/internal/utils"
	"blog/models"
	"errors"

	"go.uber.org/zap"
)

// LoginVerification 管理员登录验证 - 使用bcrypt密码比较
func LoginVerification(username, password string) (models.Admin, error) {
	var ad models.Admin
	err := core.DB.Where("username = ?", username).First(&ad).Error
	if err != nil {
		zap.L().Error("LoginVerification: 查询管理员失败 " + err.Error())
		return ad, err
	}
	// bcrypt密码校验
	if !utils.CheckPassword(ad.Password, password) {
		return models.Admin{}, errors.New("密码错误")
	}
	return ad, nil
}

// GetAdminByID 根据ID获取管理员
func GetAdminByID(id uint64) (models.Admin, error) {
	var ad models.Admin
	err := core.DB.First(&ad, id).Error
	if err != nil {
		zap.L().Error("GetAdminByID:" + err.Error())
		return ad, err
	}
	return ad, nil
}

// CreateAdmin 创建管理员（seed用，密码会bcrypt加密）
func CreateAdmin(admin models.Admin) error {
	hashedPwd, err := utils.HashPassword(admin.Password)
	if err != nil {
		zap.L().Error("CreateAdmin 密码加密失败:" + err.Error())
		return err
	}
	admin.Password = hashedPwd
	err = core.DB.Create(&admin).Error
	if err != nil {
		zap.L().Error("CreateAdmin:" + err.Error())
		return err
	}
	return nil
}
