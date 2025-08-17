package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micro-service-platform/pkg/database"
	"micro-service-platform/pkg/logger"
	"micro-service-platform/pkg/registry"
	"micro-service-platform/services/user-service/internal/config"
	"micro-service-platform/services/user-service/internal/handler"
	"micro-service-platform/services/user-service/internal/repository"
	"micro-service-platform/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化数据库
	db, err := database.NewPostgresDB(cfg.Database.DSN)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Service will start without database connection")
		db = nil
	}

	// 初始化Redis
	redisClient, err := database.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Service will start without Redis connection")
		redisClient = nil
	}

	// 初始化服务注册中心客户端
	registryClient, err := registry.NewEtcdClient(cfg.Registry.Endpoints)
	if err != nil {
		log.Printf("Warning: Failed to create registry client: %v", err)
		log.Println("Service will start without service registry")
		registryClient = nil
	} else {
		defer registryClient.Close()
	}

	// 准备服务信息
	var serviceInfo *registry.ServiceInfo
	if registryClient != nil {
		serviceInfo = &registry.ServiceInfo{
			Name:    cfg.Service.Name,
			Address: cfg.Service.Address,
			Port:    cfg.Service.Port,
			Tags:    []string{"user", "api"},
		}

		if err := registryClient.Register(context.Background(), serviceInfo); err != nil {
			log.Printf("Warning: Failed to register service: %v", err)
		} else {
			log.Println("Service registered successfully")
		}
	}

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db, redisClient)

	// 初始化服务层
	userService := service.NewUserService(userRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)

	// 创建Gin引擎
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 设置路由
	setupRoutes(r, userHandler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("User Service starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down User Service...")

	// 注销服务
	if registryClient != nil && serviceInfo != nil {
		if err := registryClient.Deregister(context.Background(), serviceInfo); err != nil {
			log.Printf("Failed to deregister service: %v", err)
		} else {
			log.Println("Service deregistered successfully")
		}
	}

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("User Service stopped")
}

func setupRoutes(r *gin.Engine, h *handler.UserHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "user-service",
		})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("/register", h.Register)
			users.POST("/login", h.Login)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
			users.GET("/", h.ListUsers)
		}
	}
}
