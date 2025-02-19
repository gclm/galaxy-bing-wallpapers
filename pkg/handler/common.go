package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	switch err {
	case mongo.ErrNoDocuments:
		c.JSON(404, ErrorResponse{
			Code:    404,
			Message: "Wallpaper not found",
		})
	default:
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(c.Writer.Status(), gin.H{
				"error":  c.Errors.Last().Err.Error(),
				"status": c.Writer.Status(),
			})
		}
	}
}
