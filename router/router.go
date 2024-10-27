package router

import (
	"order_api/app/auth"
	"order_api/handler"
	"order_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orderHandler *handler.OrderHandler, authHandler *handler.AuthHandler, authService *auth.AuthService) *gin.Engine {
	router := gin.New()

	// 添加中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 注册验证器
	handler.RegisterValidators()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 认证路由
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// API 路由
	v1 := router.Group("/api/v1")
	{
		// 添加认证中间件
		v1.Use(middleware.Auth(authService))

		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id", orderHandler.UpdateOrder)
			orders.DELETE("/:id", orderHandler.DeleteOrder)
		}
	}

	return router
}
