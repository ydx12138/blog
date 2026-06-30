package config

import (
	"blog/flags"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Cfg *Config
)

// 总的结构体
type Config struct {
	//Server       ServerConfig `mapstructure:"server"`
	CORS         CORSConfig   `mapstructure:"cors"`
	SystemConfig SystemConfig `mapstructure:"system"`
	LogConfig    LogConfig    `mapstructure:"log"`
	MysqlConfig  MysqlConfig  `mapstructure:"mysql"`
	Redis        RedisConfig  `mapstructure:"redis"`
	OssConfig    OssConfig    `mapstructure:"oss"`
	MailConfig   MailConfig   `mapstructure:"mail"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}
type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	FromName string `mapstructure:"from_name"`
	SSL      bool   `mapstructure:"ssl"`
}

type MysqlConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Db        string `mapstructure:"db"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Log_level string `mapstructure:"log_level"`
}

func (m MysqlConfig) DSN() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.Db + "?charset=utf8mb4&parseTime=true"

}

type SystemConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

func (s SystemConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type LogConfig struct {
	App   string `mapstructure:"app"`
	Dir   string `mapstructure:"dir"`
	Level string `mapstructure:"level"`
}
type OssConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	Image_path      string `mapstructure:"image_path"`
}

/*type ServerConfig struct {
	Port int `mapstructure:"port"`
}*/

type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// LoadConfig 加载配置文件
/*func LoadConfig() error {
	viper.SetConfigFile(flags.FlagOptions.File)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}
	zap.L().Info("读取配置文件" + flags.FlagOptions.File + "成功")
	return nil
}*/
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigFile(flags.FlagOptions.File)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, err
	}
	zap.L().Info("读取配置文件" + flags.FlagOptions.File + "成功")
	Cfg = cfg
	return cfg, nil
}
