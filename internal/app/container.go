package app

import (
	"blog/config"
	"blog/internal/handler"
	"blog/internal/repository"
	"blog/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Config = *config.Config

type Container struct {
	Config     Config
	DB         *gorm.DB
	Redis      *redis.Client
	Repository *repository.Repository
	Service    *service.Service
	Handler    *handler.Handler
}

func NewContainer(cfg Config, db *gorm.DB, redis *redis.Client) *Container {
	repo := repository.New(db)
	svc := service.New(repo, redis)
	h := handler.New(svc)
	return &Container{
		Config:     cfg,
		DB:         db,
		Redis:      redis,
		Repository: repo,
		Service:    svc,
		Handler:    h,
	}
}
