package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gclm/galaxy-bing-wallpapers/pkg/database"
	"github.com/gclm/galaxy-bing-wallpapers/pkg/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTodayWallpaper(c *gin.Context) {
	collection := database.GetCollection("wallpapers")

	// 处理查询参数
	width := c.DefaultQuery("w", "1920")
	height := c.DefaultQuery("h", "1080")
	mkt := c.DefaultQuery("mkt", "zh-CN")
	responseType := c.DefaultQuery("type", "image")

	var wallpaper model.Wallpaper
	ctx := context.Background()
	err := collection.FindOne(ctx, bson.M{
		"datetime": bson.M{"$gte": time.Now().Format("2006-01-02")},
		"mkt":      mkt,
	}).Decode(&wallpaper)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallpaper not found"})
		return
	}

	imageURL := wallpaper.GenerateImageURL(width, height)

	switch responseType {
	case "image":
		c.Redirect(http.StatusFound, imageURL)
	case "json":
		c.JSON(http.StatusOK, model.ImageResponse{
			Url:      imageURL,
			Title:    wallpaper.Title,
			Datetime: wallpaper.Datetime,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported response type. Use 'image' or 'json'",
		})
	}
}
