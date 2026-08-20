package config

import (
	"blog/flags"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
	Wechat       WechatConfig `mapstructure:"wechat"`
	Sms          SmsConfig    `mapstructure:"sms"`
}

type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
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
	DefaultAvatar   string `mapstructure:"default_avatar"`
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

type SmsConfig struct {
	SchemeName       string `mapstructure:"scheme_name"`        //方案名称
	CountryCode      string `mapstructure:"country_code"`       //号码国家编码
	SignName         string `mapstructure:"sign_name"`          //签名名称
	TemplateCode     string `mapstructure:"template_code"`      //短信模板CODE
	TemplateParam    string `mapstructure:"template_param"`     //短信模板
	CodeLength       int64  `mapstructure:"code_length"`        //验证码长度
	ValidTime        int64  `mapstructure:"valid_time"`         //验证码有效时间
	DuplicatePolicy  int64  `mapstructure:"duplicate_policy"`   //核验规则
	Interval         int64  `mapstructure:"interval"`           //时间间隔
	CodeType         int64  `mapstructure:"code_type"`          //生成的验证码类型
	ReturnVerifyCode bool   `mapstructure:"return_verify_code"` //是否返回验证码
	AutoRetry        int64  `mapstructure:"auto_retry"`         //是否自动替换签名重试
	AccessKeyId      string `mapstructure:"access_key_id"`
	AccessKeySecret  string `mapstructure:"access_key_secret"`
	CaseAuthPolicy   int64  `mapstructure:"case_auth_policy"`
}

func LoadConfig() (*Config, error) {

	// 1. 加载 .env 文件到环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 2. 初始化 Viper
	v := viper.New()
	v.SetConfigFile("settings.yaml")
	v.SetConfigType("yaml")

	// 3. 绑定环境变量（自动将 YAML 键名转换为环境变量）
	// 例如：server.port -> ENJOYMALL_SERVER_PORT
	v.SetEnvPrefix("BLOG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.BindEnv("mysql.db", "BLOG_MYSQL_DATABASE", "MYSQL_DATABASE"); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := v.Unmarshal(cfg); err != nil {
		return cfg, err
	}
	zap.L().Info("读取配置文件" + flags.FlagOptions.File + "成功")
	Cfg = cfg
	return cfg, nil
}

//func Load() (*Config, error) {
//	// 1. 加载 .env 文件到环境变量
//	if err := godotenv.Load(); err != nil {
//		log.Println("Warning: .env file not found, using environment variables")
//	}
//
//	// 2. 初始化 Viper
//	v := viper.New()
//	v.SetConfigName("config")
//	v.SetConfigType("yaml")
//	v.AddConfigPath("./configs")
//	v.AddConfigPath(".")
//
//	// 3. 绑定环境变量（自动将 YAML 键名转换为环境变量）
//	// 例如：server.port -> ENJOYMALL_SERVER_PORT
//	v.SetEnvPrefix("ENJOYMALL")
//	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
//	v.AutomaticEnv()
//
//	// 4. 显式绑定 Coze 配置的环境变量（处理下划线键名）
//	v.BindEnv("coze.bot_id", "ENJOYMALL_COZE_BOT_ID")
//	v.BindEnv("coze.review_summary_bot_id", "ENJOYMALL_COZE_REVIEW_SUMMARY_BOT_ID")
//	v.BindEnv("coze.access_token", "ENJOYMALL_COZE_ACCESS_TOKEN")
//	v.BindEnv("coze.api_url", "ENJOYMALL_COZE_API_URL")
//	v.BindEnv("coze.timeout", "ENJOYMALL_COZE_TIMEOUT")
//
//	// 5. 读取配置文件
//	if err := v.ReadInConfig(); err != nil {
//		log.Fatal("Error reading config.yaml:", err)
//	}
//
//	// 6. 解析配置到结构体
//
//	if err := v.Unmarshal(&cfg); err != nil {
//		log.Fatal("Error unmarshaling config:", err)
//	}
//
//	return cfg, nil
//}
