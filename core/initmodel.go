package core

import (
	model2 "blog/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitModel(db *gorm.DB) {
	err := db.AutoMigrate(
		&model2.User{},
		&model2.Admin{},
		&model2.Category{},
		&model2.Article{},
		&model2.Comment{},
		&model2.Token{},
	)
	if err != nil {
		zap.L().Panic("migrate tables failed: " + err.Error())
	}
}
