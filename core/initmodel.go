package core

import (
	model2 "blog/models"

	"go.uber.org/zap"
)

func InitModel() {
	err := DB.AutoMigrate(
		&model2.User{},
		&model2.Admin{},
		&model2.Category{},
		&model2.Article{},
		&model2.Comment{},
		&model2.Token{},
	)

	if err != nil {
		zap.L().Panic("表迁移出错" + err.Error())
	}
}
