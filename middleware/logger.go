// Package middleware 提供HTTP中间件功能
package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// bodyLogWriter 自定义响应写入器，用于记录响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 实现io.Writer接口，用于捕获响应内容
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 日志中间件，记录请求和响应信息
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始计时
		start := time.Now()

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// 恢复请求体，以供后续处理
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 创建自定义响应写入器
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算请求处理时间
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// 记录错误信息
		if statusCode >= 400 {
			gin.DefaultErrorWriter.Write([]byte(
				"[ERROR] " +
					c.Request.Method + " " +
					c.Request.URL.Path + " " +
					string(bodyBytes) + " " +
					blw.body.String() + "\n",
			))
		}

		// 记录请求摘要
		gin.DefaultWriter.Write([]byte(
			"[INFO] " +
				c.Request.Method + " " +
				c.Request.URL.Path + " " +
				duration.String() + "\n", // 使用duration.String()进行转换
		))
	}
}
