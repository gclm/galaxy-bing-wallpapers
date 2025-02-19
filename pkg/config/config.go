package config

import (
	"fmt"
	"os"
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
func LoadConfig() (*Config, error) {
	var err error

	// 使用 sync.Once 确保全局配置只被初始化一次
	once.Do(func() {
		// 仅在本地开发环境加载 .env 文件
		if os.Getenv("VERCEL") != "1" {
			if err = godotenv.Load(); err != nil {
				fmt.Println("Warning: .env file not found, using environment variables")
			}
		}

		// 获取并验证 MongoDB URI
		mongoDBURI := os.Getenv("MONGODB_URI")
		if mongoDBURI == "" {
			err = fmt.Errorf("MONGODB_URI environment variable is required")
			return
		}

		GlobalConfig = &Config{
			MongoDBURI: mongoDBURI,
			APIToken:   os.Getenv("API_TOKEN"),
			GinMode:    getEnvWithDefault("GIN_MODE", "release"),
		}
	})

	if err != nil {
		return nil, err
	}

	return GlobalConfig, nil
}

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
