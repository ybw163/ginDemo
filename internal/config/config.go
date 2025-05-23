package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// 服务器配置
	Server ServerConfig `destructure:"server"`
	//数据库配置
	Database DatabaseConfig `destructure:"database"`
	// JWT配置
	JWT JWTConfig `destructure:"jwt"`
	// 日志配置
	Log LogConfig `destructure:"log"`
	// Redis配置
	Redis RedisConfig `destructure:"redis"`
}

type ServerConfig struct {
	Port         string `destructure:"port"`
	Mode         string `destructure:"mode"` // debug, release, test
	ReadTimeout  int    `destructure:"read_timeout"`
	WriteTimeout int    `destructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver   string `destructure:"driver"`
	Host     string `destructure:"host"`
	Port     string `destructure:"port"`
	Username string `destructure:"username"`
	Password string `destructure:"password"`
	Database string `destructure:"database"`
	Charset  string `destructure:"charset"`
}

type JWTConfig struct {
	Secret     string `destructure:"secret"`
	ExpireTime int    `destructure:"expire_time"` // hours
}

type LogConfig struct {
	Level string `destructure:"level"`
	Path  string `destructure:"path"`
}

type RedisConfig struct {
	Host     string `destructure:"host"`
	Port     string `destructure:"port"`
	Password string `destructure:"password"`
	DB       int    `destructure:"db"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("jwt.expire_time", 24)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.path", "./logs")
}
