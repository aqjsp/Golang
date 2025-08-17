package router

import (
	"micro-service-platform/api-gateway/internal/handler"
	"micro-service-platform/api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, h *handler.Handler) {
	// 健康检查
	r.GET("/health", h.Health)

	// 指标端点
	r.GET("/metrics", h.Metrics)

	// 服务发现端点
	r.GET("/services", h.GetServices)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户服务路由
		userGroup := api.Group("/users")
		{
			userGroup.Use(middleware.Metrics())
			userGroup.Any("/*path", h.ProxyRequest("user-service"))
		}

		// 订单服务路由
		orderGroup := api.Group("/orders")
		{
			orderGroup.Use(middleware.Auth()) // 订单需要认证
			orderGroup.Use(middleware.Metrics())
			orderGroup.Any("/*path", h.ProxyRequest("order-service"))
		}

		// 支付服务路由
		paymentGroup := api.Group("/payments")
		{
			paymentGroup.Use(middleware.Auth()) // 支付需要认证
			paymentGroup.Use(middleware.Metrics())
			paymentGroup.Any("/*path", h.ProxyRequest("payment-service"))
		}
	}

	// 管理API路由组
	admin := r.Group("/admin")
	{
		admin.Use(middleware.Auth()) // 管理接口需要认证
		admin.GET("/services", h.GetServices)
		admin.GET("/health", h.Health)
	}

	// WebSocket路由（用于实时通信）
	r.GET("/ws", handleWebSocket)

	// 静态文件服务（如果需要）
	r.Static("/static", "./static")

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})
}

// handleWebSocket WebSocket处理器
func handleWebSocket(c *gin.Context) {
	// TODO: 实现WebSocket处理逻辑
	c.JSON(200, gin.H{
		"message": "WebSocket endpoint - TODO: implement WebSocket handler",
	})
}
