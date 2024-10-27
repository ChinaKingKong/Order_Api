// Package middleware 提供HTTP中间件功能
package middleware

import (
	"net/http"
	"order_api/app/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件，用于验证请求中的JWT令牌
func Auth(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "缺少Authorization头",
			})
			return
		}

		// 解析Bearer令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效的Authorization头格式",
			})
			return
		}

		// 验证令牌
		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效或已过期的令牌",
			})
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
