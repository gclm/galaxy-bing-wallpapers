package handler

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gclm/galaxy-bing-api/pkg/database"
	"github.com/gclm/galaxy-bing-api/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRandomWallpaper(c *gin.Context) {
	collection := database.GetCollection("wallpapers")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取总数量
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 生成随机索引
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(int(total))

	// 获取随机文档
	var result bson.M
	err = collection.FindOne(ctx, bson.M{}, options.FindOne().SetSkip(int64(randomIndex))).Decode(&result)

	if handleMongoError(c, err) {
		return
	}

	redirectToImage(c, result)
}

func handleMongoError(c *gin.Context, err error) bool {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallpaper not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return true
	}
	return false
}

func redirectToImage(c *gin.Context, result bson.M) {
	width := c.DefaultQuery("w", "1920")
	height := c.DefaultQuery("h", "1080")

	var wallpaper model.Wallpaper
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &wallpaper)

	imageURL := wallpaper.GenerateImageURL(width, height)

	responseType := c.DefaultQuery("type", "image")

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
