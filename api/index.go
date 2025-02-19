package handler

import (
	"net/http"

	"github.com/gclm/galaxy-bing-api/internal/config"
	"github.com/gclm/galaxy-bing-api/internal/database"
	"github.com/gclm/galaxy-bing-api/internal/handler"
	"github.com/gclm/galaxy-bing-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 初始化 Gin
	gin.SetMode(cfg.GinMode)
	app = gin.New()
	app.Use(middleware.CorsMiddleware())
	app.Use(middleware.Recovery())

	// 初始化数据库
	if err := database.InitMongoDB(); err != nil {
		panic(err)
	}

	// 注册路由
	setupRoutes(app)
}

// Handler Vercel serverless function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
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
