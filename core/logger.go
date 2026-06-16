package core

import (
	"blog/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// logEncoder 时间分片和level分片同时做
type logEncoder struct {
	zapcore.Encoder
	errFile     *os.File
	file        *os.File
	currentDate string
}

// 这里相当于是截住了日志，做出手脚
func (e *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行，加前缀
	buff, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	data := buff.String()
	buff.Reset()
	buff.AppendString("[" + config.Cfg.LogConfig.App + "] " + data)
	data = buff.String()
	// 时间分片
	now := time.Now().Format("2006-01-02")
	if e.currentDate != now {
		//
		if e.file != nil {
			err = e.file.Close()
			if err != nil {
				zap.L().Error("close file error", zap.Error(err))
			}
		}
		//检查目录

		err = os.MkdirAll(fmt.Sprintf(config.Cfg.LogConfig.Dir+"/%s", now), 0666)
		if err != nil {
			return nil, err
		}
		// 时间不同，先创建目录，并打开文件
		name := fmt.Sprintf(config.Cfg.LogConfig.Dir+"/%s/out.log", now)
		var file *os.File
		file, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		e.file = file
		e.currentDate = now
	}
	//把err存到err.log

	switch entry.Level {
	case zapcore.ErrorLevel:
		if e.errFile == nil {
			name := fmt.Sprintf(config.Cfg.LogConfig.Dir+"/%s/err.log", now)
			file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			e.errFile = file
		}
		_, err = e.errFile.WriteString(data)
		if err != nil {
			return nil, err
		}
	default:
	}
	//out.log里存全部
	if e.currentDate == now {
		_, err2 := e.file.WriteString(data)
		if err2 != nil {
			return nil, err2
		}
	}
	return buff, nil
}

func LogInit() *zap.Logger {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	//cfg.EncoderConfig.EncodeLevel = myEncodeLevel
	// 创建自定义的 Encoder
	encoder := &logEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器，自定义分片
	}
	// 创建 Core
	var level zapcore.Level = zapcore.DebugLevel
	switch config.Cfg.LogConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}
	fmt.Println(level)
	core := zapcore.NewCore(
		encoder,                                // 使用自定义的 Encoder
		zapcore.NewMultiWriteSyncer(os.Stdout), // 输出到控制台
		level,                                  // 设置日志级别
	)
	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
	return logger
}
