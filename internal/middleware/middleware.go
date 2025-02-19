package middleware

import (
	"github.com/gclm/galaxy-bing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(500, handler.ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		})
	})
}
