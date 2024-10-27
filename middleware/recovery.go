// Package middleware 提供HTTP中间件功能
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件，用于捕获和处理panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息并返回500错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}