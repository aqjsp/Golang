package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"micro-service-platform/api-gateway/internal/config"
	"micro-service-platform/pkg/registry"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	registryClient registry.Client
	config         *config.Config
	httpClient     *http.Client
}

func NewHandler(registryClient registry.Client, cfg *config.Config) *Handler {
	return &Handler{
		registryClient: registryClient,
		config:         cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "api-gateway",
	})
}

// ProxyRequest 代理请求到后端服务
func (h *Handler) ProxyRequest(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从服务注册中心获取服务实例
		instances, err := h.registryClient.Discover(context.Background(), serviceName)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to discover service: %s", serviceName)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service unavailable",
			})
			return
		}

		if len(instances) == 0 {
			logrus.Errorf("No instances found for service: %s", serviceName)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "No service instances available",
			})
			return
		}

		// 简单的轮询负载均衡
		instance := instances[0] // 这里可以实现更复杂的负载均衡算法

		// 构建目标URL
		targetURL := fmt.Sprintf("http://%s%s", instance.Address, c.Request.URL.Path)
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// 转发请求
		h.forwardRequest(c, targetURL)
	}
}

// forwardRequest 转发HTTP请求
func (h *Handler) forwardRequest(c *gin.Context, targetURL string) {
	// 读取请求体
	var body []byte
	if c.Request.Body != nil {
		body, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// 创建新的请求
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
	if err != nil {
		logrus.WithError(err).Error("Failed to create request")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 添加追踪头
	req.Header.Set("X-Request-ID", generateRequestID())
	req.Header.Set("X-Forwarded-For", c.ClientIP())

	// 发送请求
	resp, err := h.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to forward request")
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Bad gateway",
		})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	io.Copy(c.Writer, resp.Body)
}

// GetServices 获取所有注册的服务
func (h *Handler) GetServices(c *gin.Context) {
	services := []string{"user-service", "order-service", "payment-service"}
	result := make(map[string]interface{})

	for _, serviceName := range services {
		instances, err := h.registryClient.Discover(context.Background(), serviceName)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to discover service: %s", serviceName)
			continue
		}
		result[serviceName] = instances
	}

	c.JSON(http.StatusOK, gin.H{
		"services": result,
	})
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Metrics 获取网关指标
func (h *Handler) Metrics(c *gin.Context) {
	// TODO: 实现Prometheus指标输出
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics endpoint - TODO: implement Prometheus metrics",
	})
}