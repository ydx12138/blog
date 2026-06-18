package dao

import (
	"blog/core"
	"blog/models"

	"go.uber.org/zap"
)

func LoginVerification(username, password string) (models.Admin, error) {
	var ad models.Admin
	err := core.DB.Find(&ad, "username = ? AND password = ?", username, password).Error
	if err != nil {
		zap.L().Error("LoginVerification:" + err.Error())
		return ad, err
	}
	return ad, nil
}
