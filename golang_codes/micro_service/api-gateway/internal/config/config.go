package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Registry  RegistryConfig  `mapstructure:"registry"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Log       LogConfig       `mapstructure:"log"`
	Services  ServicesConfig  `mapstructure:"services"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type RegistryConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type ServicesConfig struct {
	UserService    ServiceConfig `mapstructure:"user_service"`
	OrderService   ServiceConfig `mapstructure:"order_service"`
	PaymentService ServiceConfig `mapstructure:"payment_service"`
}

type ServiceConfig struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}

func setDefaults() {
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("registry.endpoints", []string{"localhost:2379"})
	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("services.user_service.name", "user-service")
	viper.SetDefault("services.user_service.path", "/api/v1/users")
	viper.SetDefault("services.order_service.name", "order-service")
	viper.SetDefault("services.order_service.path", "/api/v1/orders")
	viper.SetDefault("services.payment_service.name", "payment-service")
	viper.SetDefault("services.payment_service.path", "/api/v1/payments")
}
