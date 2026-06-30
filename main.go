package main

import (
	"blog/config"
	"blog/core"
	"blog/flags"
	"blog/internal/app"
	"blog/internal/router"
	"blog/seed"

	"go.uber.org/zap"
)

func main() {
	flags.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Error("config load failed: " + err.Error())
		return
	}

	core.LogInit()

	db, err := core.DataBaseInit()
	if err != nil {
		return
	}

	core.InitModel(db)

	if flags.FlagOptions.Seed {
		seed.Run()
		zap.L().Info("seed data completed")
		return
	}

	container := app.NewContainer(cfg, db, nil)
	r := router.Register(container.Handler)

	zap.L().Debug("gin running at " + config.Cfg.SystemConfig.Address())
	if err := r.Run(config.Cfg.SystemConfig.Address()); err != nil {
		zap.L().Error("router run failed: " + err.Error())
	}
}
