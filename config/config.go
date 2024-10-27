package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 系统配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Log      LogConfig      `json:"log"`
	JWT      JWTConfig      `json:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            string `json:"port"`
	ReadTimeout     int    `json:"read_timeout"`
	WriteTimeout    int    `json:"write_timeout"`
	ShutdownTimeout int    `json:"shutdown_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DBName       string `json:"dbname"`
	Charset      string `json:"charset"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	MaxRetries   int    `json:"max_retries"`
	PoolSize     int    `json:"pool_size"`
	MaxIdleConns int    `json:"max_idle_conns"`
	ExpireHours  int    `json:"expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey          string `json:"secret_key"`
	TokenExpiryHours   int    `json:"token_expiry_hours"`
	RefreshExpiryHours int    `json:"refresh_expiry_hours"`
}

// NewConfig 创建新的配置实例
func NewConfig() *Config {
	config := &Config{}
	if err := config.Load(); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}
	return config
}

// Load 加载配置文件
func (c *Config) Load() error {
	// 读取配置文件
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析JSON配置
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := c.validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// validate 验证配置
func (c *Config) validate() error {
	// 验证服务器配置
	if c.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	// 验证数据库配置
	if c.Database.Host == "" || c.Database.Port == "" ||
		c.Database.User == "" || c.Database.DBName == "" {
		return fmt.Errorf("数据库配置不完整")
	}

	// 验证Redis配置
	if c.Redis.Host == "" || c.Redis.Port == "" {
		return fmt.Errorf("Redis配置不完整")
	}

	// 验证JWT配置
	if c.JWT.SecretKey == "" || c.JWT.TokenExpiryHours <= 0 {
		return fmt.Errorf("JWT配置不完整")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
	)
}

// GetRedisAddr 获取Redis连接地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
