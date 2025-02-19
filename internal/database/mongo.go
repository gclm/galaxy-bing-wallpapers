package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gclm/galaxy-bing-api/internal/config"
	"github.com/gclm/galaxy-bing-api/internal/model"
)

var Client *mongo.Client

func InitMongoDB() error {
	uri := config.GlobalConfig.MongoDBURI
	if uri == "" {
		return fmt.Errorf("MONGODB_URI environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	Client = client
	return nil
}

func GetCollection(name string) *mongo.Collection {
	return Client.Database("bing").Collection(name)
}

// SaveWallpaper 保存壁纸信息到数据库
func SaveWallpaper(wallpaper model.Wallpaper) error {
	collection := GetCollection("wallpapers")
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
			options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})).Decode(&lastWallpaper)

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

		log.Printf("Inserted new wallpaper: ID=%d, Title=%s", wallpaper.ID, wallpaper.Title)
	}

	return nil
}

// CreateIndexes 创建必要的数据库索引
func CreateIndexes() error {
	collection := GetCollection("wallpapers")
	ctx := context.Background()

	// 创建ID唯一索引
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return fmt.Errorf("failed to create id index: %v", err)
	}

	// 创建日期和市场代码复合索引
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "datetime", Value: 1},
			{Key: "mkt", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return fmt.Errorf("failed to create datetime-mkt index: %v", err)
	}

	return nil
}
