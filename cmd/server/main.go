package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gclm/galaxy-bing-api/pkg/config"
	"github.com/gclm/galaxy-bing-api/pkg/database"
	"github.com/gclm/galaxy-bing-api/pkg/handler"
	"github.com/gclm/galaxy-bing-api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.GinMode)

	// 使用环境变量中的端口，默认为 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 设置工作目录为项目根目录
	if err := os.Chdir(projectRoot()); err != nil {
		panic(fmt.Sprintf("Failed to change working directory: %v", err))
	}

	// 初始化数据库
	if err := database.InitMongoDB(); err != nil {
		panic(err)
	}

	// 初始化 Gin
	app := gin.New()
	app.Use(middleware.CorsMiddleware())
	app.Use(middleware.Recovery())

	// 注册路由
	setupRoutes(app)

	// 启动服务
	addr := fmt.Sprintf(":%s", port)
	if err := app.Run(addr); err != nil {
		panic(err)
	}
}

// setupRoutes 设置路由
func setupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/today", handler.GetTodayWallpaper)
		v1.GET("/random", handler.GetRandomWallpaper)
		v1.GET("/list", middleware.TokenAuth(), handler.GetWallpaperList)
		v1.GET("/date/:date", middleware.TokenAuth(), handler.GetWallpaperByDate)
		v1.GET("/health", handler.HealthCheck)
	}
}

// projectRoot 获取项目根目录
func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}
