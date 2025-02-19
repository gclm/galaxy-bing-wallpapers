package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gclm/galaxy-bing-api/internal/database"
	"github.com/gclm/galaxy-bing-api/internal/model"
)

func main() {
	// 初始化数据库连接
	if err := database.InitMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 创建数据库索引
	if err := database.CreateIndexes(); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// 读取data目录
	dataDir := "data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
	}

	// 遍历所有JSON文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_all.json") {
			mkt := strings.TrimSuffix(file.Name(), "_all.json")
			log.Printf("Processing %s market data...", mkt)

			// 读取文件内容
			filePath := filepath.Join(dataDir, file.Name())
			if err := importDataFile(filePath, mkt); err != nil {
				log.Printf("Failed to import data for %s: %v", mkt, err)
				continue
			}

			log.Printf("Successfully imported data for %s market", mkt)
		}
	}
}

func importDataFile(filePath, mkt string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 解析JSON数据
	var response model.WallpaperList
	if err := json.Unmarshal(content, &response); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// 批量导入数据
	for _, wallpaper := range response.Data {
		// 确保设置了市场代码
		wallpaper.Mkt = mkt

		// 保存到数据库
		if err := database.SaveWallpaper(wallpaper); err != nil {
			log.Printf("Warning: failed to save wallpaper %s: %v", wallpaper.Title, err)
			continue
		}
	}

	log.Printf("Successfully imported %d wallpapers for %s market", len(response.Data), mkt)

	return nil
}
