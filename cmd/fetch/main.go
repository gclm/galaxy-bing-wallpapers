package main

import (
	"log"

	"github.com/gclm/galaxy-bing-wallpapers/pkg/config"
	"github.com/gclm/galaxy-bing-wallpapers/pkg/database"
	"github.com/gclm/galaxy-bing-wallpapers/pkg/utils"
)

func main() {
	// 加载配置
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 支持的市场代码
	markets := []string{
		"zh-CN", // 中国
		"de-DE", // 德国
		"en-CA", // 加拿大（英语）
		"en-GB", // 英国
		"en-IN", // 印度
		"en-US", // 美国
		"fr-FR", // 法国
		"it-IT", // 意大利
		"ja-JP", // 日本
	}

	// 获取每个市场的壁纸
	for _, mkt := range markets {
		log.Printf("正在获取 %s 市场的壁纸...", mkt)
		isNew, err := utils.FetchLatestWallpaper(mkt)
		if err != nil {
			log.Printf("获取 %s 市场壁纸失败: %v", mkt, err)
			continue
		}

		if isNew {
			log.Printf("✅ %s 市场壁纸已成功保存到数据库", mkt)
		} else {
			log.Printf("ℹ️ %s 市场今日壁纸已存在，跳过保存", mkt)
		}
	}
}
