package main

import (
	"blog/config"
	"blog/core"
	"blog/flags"
	"blog/router"

	"go.uber.org/zap"
)

// 先把many to many去掉
func main() {
	//参数
	flags.Parse()
	//加载配置文件
	err := config.LoadConfig()
	if err != nil {
		zap.L().Error("配置文件加载失败" + err.Error())
		return
	}
	//加载日志
	core.LogInit()
	//连接数据库
	core.DataBaseInit()
	//迁移表
	//core.InitModel()
	//注册路由
	r := router.Register()
	zap.L().Debug("gin运行在" + config.Cfg.SystemConfig.Address())
	err = r.Run(config.Cfg.SystemConfig.Address())
	if err != nil {
		zap.L().Error("路由加载失败" + err.Error())
		return
	}
}
