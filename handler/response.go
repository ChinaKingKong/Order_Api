// Package handler 提供HTTP请求处理器
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data,omitempty"`    // 数据
	Errors  []string    `json:"errors,omitempty"`  // 错误信息
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
	})
}

// Created 返回创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string, errors ...string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}

// ValidationError 返回验证错误响应
func ValidationError(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "数据验证失败",
		Errors:  errors,
	})
}

// ServerError 返回服务器错误响应
func ServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
		Errors:  []string{err.Error()},
	})
}

// Unauthorized 返回未授权响应
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问",
	})
}

// Forbidden 返回禁止访问响应
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: "禁止访问",
	})
}

// NotFound 返回资源未找到响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}