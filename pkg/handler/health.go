package handler

import (
	"github.com/gclm/galaxy-bing-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	// 检查数据库连接
	if err := database.Client.Ping(c, nil); err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: "Database connection error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}
