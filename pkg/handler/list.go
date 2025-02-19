package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gclm/galaxy-bing-wallpapers/pkg/database"
	"github.com/gclm/galaxy-bing-wallpapers/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllWallpapers(c *gin.Context) {
	collection := database.GetCollection("wallpapers")
	ctx := context.Background()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "startdate", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if handleMongoError(c, err) {
		return
	}

	var results []model.Wallpaper
	if err = cursor.All(ctx, &results); handleMongoError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"pagination": gin.H{
			"currentPage": page,
			"pageSize":    pageSize,
			"total":       getTotalCount(collection, ctx),
		},
	})
}

func GetTotalCount(c *gin.Context) {
	collection := database.GetCollection("wallpapers")
	ctx := context.Background()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if handleMongoError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count})
}

func getTotalCount(collection *mongo.Collection, ctx context.Context) int64 {
	count, _ := collection.CountDocuments(ctx, bson.M{})
	return count
}

// GetWallpaperList 获取壁纸列表
func GetWallpaperList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	mkt := c.DefaultQuery("mkt", "")

	// 转换为整数
	skip, limit := getPagination(page, pageSize)

	// 构建查询条件
	filter := bson.M{}
	if mkt != "" {
		filter["mkt"] = mkt
	}

	// 获取总数
	total, wallpapers := getWallpapers(c, filter, skip, limit)

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    wallpapers,
		Total:   total,
	})
}

// GetWallpaperByDate 获取指定日期的壁纸
func GetWallpaperByDate(c *gin.Context) {
	date := c.Param("date") // 格式：2024-02-19
	mkt := c.Query("mkt")   // 可选参数
	width := c.DefaultQuery("w", "1920")
	height := c.DefaultQuery("h", "1080")
	responseType := c.DefaultQuery("type", "image")

	collection := database.GetCollection("wallpapers")
	ctx := context.Background()

	// 构建查询条件
	filter := bson.M{"datetime": date}
	if mkt != "" {
		filter["mkt"] = mkt
	}

	// 查询壁纸
	var wallpaper model.Wallpaper
	err := collection.FindOne(ctx, filter).Decode(&wallpaper)
	if err != nil {
		HandleError(c, err)
		return
	}

	// 生成图片URL
	imageURL := wallpaper.GenerateImageURL(width, height)

	// 根据响应类型返回数据
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
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unsupported response type. Use 'image' or 'json'",
		})
	}
}

// 辅助函数：获取分页参数
func getPagination(page, pageSize string) (int64, int64) {
	p, _ := strconv.ParseInt(page, 10, 64)
	ps, _ := strconv.ParseInt(pageSize, 10, 64)
	if p < 1 {
		p = 1
	}
	if ps < 1 {
		ps = 20
	}
	return (p - 1) * ps, ps
}

// 辅助函数：获取壁纸列表
func getWallpapers(c *gin.Context, filter bson.M, skip, limit int64) (int64, []model.Wallpaper) {
	collection := database.GetCollection("wallpapers")
	ctx := context.Background()

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get total count",
		})
		return 0, nil
	}

	// 查询数据
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "datetime", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to query wallpapers",
		})
		return 0, nil
	}
	defer cursor.Close(ctx)

	var wallpapers []model.Wallpaper
	if err = cursor.All(ctx, &wallpapers); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to decode wallpapers",
		})
		return 0, nil
	}

	return total, wallpapers
}
