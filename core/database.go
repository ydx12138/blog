package core

import (
	"blog/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 结构体字段名需要大写
var DB *gorm.DB

func DataBaseInit() {
	//数据库
	//DSN := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8&parseTime=true"
	//DSN := config.Cfg.MysqlConfig.User + ":" + config.Cfg.MysqlConfig.Password + "@tcp(" + config.Cfg.MysqlConfig.Host + ":" + strconv.Itoa(config.Cfg.MysqlConfig.Port) + ")/" + config.Cfg.MysqlConfig.Db + "?charset=utf8&parseTime=true"
	var err error
	var level logger.LogLevel = logger.Info
	switch config.Cfg.MysqlConfig.Log_level {
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	case "silent":
		level = logger.Silent
	}
	DB, err = gorm.Open(mysql.Open(config.Cfg.MysqlConfig.DSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true, //不创建外键约束
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.L().Error("数据库连接失败" + err.Error())
		return
	}
}
