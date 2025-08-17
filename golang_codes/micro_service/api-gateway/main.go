package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micro-service-platform/api-gateway/internal/config"
	"micro-service-platform/api-gateway/internal/handler"
	"micro-service-platform/api-gateway/internal/middleware"
	"micro-service-platform/api-gateway/internal/router"
	"micro-service-platform/pkg/logger"
	"micro-service-platform/pkg/registry"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化服务注册中心客户端
	registryClient, err := registry.NewEtcdClient(cfg.Registry.Endpoints)
	if err != nil {
		log.Printf("Warning: Failed to create registry client: %v", err)
		log.Println("API Gateway will start without service registry")
		registryClient = nil
	} else {
		defer registryClient.Close()
	}

	// 创建Gin引擎
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(cfg.RateLimit.RequestsPerSecond))

	// 创建处理器
	h := handler.NewHandler(registryClient, cfg)

	// 设置路由
	router.SetupRoutes(r, h)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("API Gateway starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway stopped")
}
