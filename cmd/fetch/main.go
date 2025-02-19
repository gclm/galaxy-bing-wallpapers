package main

import (
	"log"

	"github.com/gclm/galaxy-bing-api/pkg/database"
	"github.com/gclm/galaxy-bing-api/pkg/utils"
)

func main() {
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
		log.Printf("Fetching wallpaper for market: %s", mkt)
		if err := utils.FetchLatestWallpaper(mkt); err != nil {
			log.Printf("Failed to fetch wallpaper for %s: %v", mkt, err)
			continue
		}
		log.Printf("Successfully fetched wallpaper for market: %s", mkt)
	}
}
