package middleware

import (
	"github.com/gclm/galaxy-bing-api/pkg/config"
	"github.com/gclm/galaxy-bing-api/pkg/handler"
	"github.com/gin-gonic/gin"
)

// TokenAuth Token 认证中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, handler.ErrorResponse{
				Code:    401,
				Message: "Authorization token is required",
			})
			c.Abort()
			return
		}

		if token != config.GlobalConfig.APIToken {
			c.JSON(403, handler.ErrorResponse{
				Code:    403,
				Message: "Invalid authorization token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
