package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	MongoDBURI string
	APIToken   string
	GinMode    string
}

// GlobalConfig 全局配置实例
var (
	GlobalConfig *Config
	once         sync.Once
)

// LoadConfig 加载并验证配置
// 配置加载优先级：
// 1. 环境变量
// 2. .env 文件
// 3. 默认值
func LoadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		// 尝试加载 .env 文件，但不强制要求
		loadEnvFile()

		// 初始化配置
		GlobalConfig = &Config{
			MongoDBURI: getRequiredEnv("MONGODB_URI"),
			APIToken:   getEnvWithDefault("API_TOKEN", "FuO2wOA4d6KUYvry"),
			GinMode:    getEnvWithDefault("GIN_MODE", "release"),
		}

		// 验证必需的配置
		if GlobalConfig.MongoDBURI == "" {
			err = fmt.Errorf("MONGODB_URI is required but not set")
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return GlobalConfig, nil
}

// loadEnvFile 尝试加载 .env 文件
// 按以下顺序查找 .env 文件：
// 1. 当前目录
// 2. 项目根目录
func loadEnvFile() {
	// 如果环境变量已存在，跳过 .env 加载
	if os.Getenv("MONGODB_URI") != "" {
		return
	}

	// 尝试加载当前目录的 .env
	if err := godotenv.Load(); err == nil {
		return
	}

	// 尝试加载项目根目录的 .env
	if root := findProjectRoot(); root != "" {
		envPath := filepath.Join(root, ".env")
		_ = godotenv.Load(envPath)
	}
}

// findProjectRoot 查找项目根目录
// 通过查找 go.mod 文件来确定项目根目录
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// getRequiredEnv 获取必需的环境变量
func getRequiredEnv(key string) string {
	return os.Getenv(key)
}

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
