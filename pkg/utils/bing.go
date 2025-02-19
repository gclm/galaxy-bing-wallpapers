package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gclm/galaxy-bing-api/pkg/database"
	"github.com/gclm/galaxy-bing-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	bingAPIURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=%s"
)

type BingResponse struct {
	Images []struct {
		URL          string `json:"url"`
		Title        string `json:"title"`
		Copyright    string `json:"copyright"`
		CopyrightURL string `json:"copyrightlink"`
		StartDate    string `json:"startdate"`
		Hsh          string `json:"hsh"`
	} `json:"images"`
}

// FetchLatestWallpaper 获取最新壁纸
func FetchLatestWallpaper(mkt string) error {
	// 构建请求URL
	url := fmt.Sprintf(bingAPIURL, mkt)

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Bing API: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析JSON响应
	var bingResp BingResponse
	if err := json.Unmarshal(body, &bingResp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(bingResp.Images) == 0 {
		return fmt.Errorf("no images found in response")
	}

	// 获取最新图片信息
	image := bingResp.Images[0]

	// 构建壁纸对象
	wallpaper := model.Wallpaper{
		Title:         image.Title,
		Url:           "https://www.bing.com" + image.URL,
		Datetime:      time.Now().Format("2006-01-02"),
		Copyright:     image.Copyright,
		CopyrightLink: image.CopyrightURL,
		Hsh:           image.Hsh,
		CreatedTime:   time.Now().Format("2006-01-02"),
		Mkt:           mkt,
	}

	// 保存到数据库
	return SaveWallpaper(wallpaper)
}

// SaveWallpaper 保存壁纸信息到数据库
func SaveWallpaper(wallpaper model.Wallpaper) error {
	collection := database.GetCollection("wallpapers")
	ctx := context.Background()

	// 检查是否已存在
	var existing model.Wallpaper
	err := collection.FindOne(ctx, bson.M{
		"datetime": wallpaper.Datetime,
		"mkt":      wallpaper.Mkt,
	}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// 获取最大ID
		var lastWallpaper model.Wallpaper
		err = collection.FindOne(ctx, bson.M{},
			options.FindOne().SetSort(bson.M{"id": -1})).Decode(&lastWallpaper)

		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("failed to get last wallpaper: %v", err)
		}

		// 设置新ID
		wallpaper.ID = lastWallpaper.ID + 1
		if wallpaper.ID == 0 {
			wallpaper.ID = 1
		}

		// 插入新记录
		_, err = collection.InsertOne(ctx, wallpaper)
		if err != nil {
			return fmt.Errorf("failed to insert wallpaper: %v", err)
		}
	}

	return nil
}
