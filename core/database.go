package core

import (
	"blog/config"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*var (
	DB *gorm.DB
)

type GormZapLogger struct {
	ZapLogger     *zap.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// NewGormZapLogger 创建 GORM 的 Zap 适配器
// level: logger.Info / Warn / Error / Silent
// slowThreshold: 慢查询阈值，例如 200*time.Millisecond
// New函数
func NewGormZapLogger(zapLogger *zap.Logger, level logger.LogLevel, slowThreshold time.Duration) *GormZapLogger {
	return &GormZapLogger{
		ZapLogger:     zapLogger,
		LogLevel:      level,
		SlowThreshold: slowThreshold,
	}
}

// LogMode 实现动态调整日志级别
func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 普通信息日志
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

// Warn 警告日志（如慢查询）
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

// Error 错误日志（如 SQL 执行错误）
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

// Trace 记录 SQL 执行细节（这是 GORM 最核心的日志埋点）
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	//开始计时
	elapsed := time.Since(begin)
	//sql语句和影响的行数
	sql, rows := fc()
	fields := []zap.Field{
		zap.Duration("latency", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && err != gorm.ErrRecordNotFound:
		l.ZapLogger.Error("SQL execution error", append(fields, zap.Error(err))...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.ZapLogger.Warn("Slow SQL query", fields...)
	case l.LogLevel >= logger.Info:
		l.ZapLogger.Info("SQL query", fields...)
	}
}

func DataBaseInit() {
	// 确保 Zap 已经初始化（你已通过 LogInit 设置了全局 logger）
	zapLogger := zap.L() // 使用全局 logger

	// 配置 GORM 的 Zap 适配器
	// 设置 Info 级别会记录所有 SQL；如果需要安静一些可以改为 Warn 或 Error
	gormLogger := NewGormZapLogger(zapLogger, logger.Info, 200*time.Millisecond)

	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8&parseTime=true"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, // 替换为 Zap 适配器
	})
	if err != nil {
		zap.L().Error("GORM 数据库连接失败", zap.Error(err))
		return
	}
}*/

// 结构体字段名需要大写
var DB *gorm.DB

func DataBaseInit() {
	//数据库
	//DSN := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8&parseTime=true"
	DSN := config.Cfg.MysqlConfig.User + ":" + config.Cfg.MysqlConfig.Password + "@tcp(" + config.Cfg.MysqlConfig.Host + ":" + strconv.Itoa(config.Cfg.MysqlConfig.Port) + ")/" + config.Cfg.MysqlConfig.Db + "?charset=utf8&parseTime=true"
	var err error
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.L().Error("数据库连接失败" + err.Error())
		return
	}
}

/*// DateFileWriter 按日期创建目录并写入 sql.log (实现 io.Writer)
type DateFileWriter struct {
	mu      sync.Mutex
	dir     string
	curDate string
	file    *os.File
}

func NewDateFileWriter(rootDir string) *DateFileWriter {
	return &DateFileWriter{dir: rootDir}
}

func (w *DateFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	date := now.Format("2006-01-02")

	if date != w.curDate {
		if w.file != nil {
			w.file.Close()
		}
		subDir := filepath.Join(w.dir, date)
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return 0, fmt.Errorf("创建日志目录失败: %w", err)
		}
		logPath := filepath.Join(subDir, "sql.log")
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, fmt.Errorf("打开日志文件失败: %w", err)
		}
		w.file = file
		w.curDate = date
	}
	return w.file.Write(p)
}

// writerAdapter 将 io.Writer 适配为 logger.Writer (实现 Printf)
type writerAdapter struct {
	io.Writer
}

func (w writerAdapter) Printf(format string, args ...interface{}) {
	// GORM 传过来的格式是类似 "[info] %s\n" 的字符串
	// 直接格式化后写入底层的 io.Writer
	fmt.Fprintf(w.Writer, format, args...)
}

var DB *gorm.DB

func DataBaseInit() {
	// 创建同时写控制台和文件的 io.Writer
	dateWriter := NewDateFileWriter("log")
	multiWriter := io.MultiWriter(os.Stdout, dateWriter)

	// 将 io.Writer 适配为 logger.Writer
	adapter := writerAdapter{Writer: multiWriter}

	// 配置 GORM 日志
	newLogger := logger.New(
		adapter, // 现在传入的是实现了 Printf 的类型
		logger.Config{
			Colorful: false,       // 控制台彩色输出（不影响文件）
			LogLevel: logger.Warn, // 可根据需要调整
		},
	)

	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8&parseTime=true"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
}*/
