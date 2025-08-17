# 分布式微服务架构平台

一个基于Go语言的企业级分布式微服务架构平台，提供完整的服务治理、监控、配置管理等功能。

## 🏗️ 项目架构

```
micro_service/
├── api-gateway/              # API网关
├── services/                 # 业务服务
│   ├── user-service/         # 用户服务
│   ├── order-service/        # 订单服务
│   └── payment-service/      # 支付服务
├── infrastructure/           # 基础设施
│   ├── service-registry/     # 服务注册中心
│   ├── config-center/        # 配置中心
│   ├── monitoring/           # 监控系统
│   └── tracing/             # 链路追踪
├── pkg/                     # 公共库
├── deployments/             # 部署配置
├── scripts/                 # 脚本文件
└── docs/                    # 文档
```

## 🚀 核心功能

- **服务注册与发现**: 基于etcd的服务注册中心
- **API网关**: 统一入口，支持路由、限流、熔断
- **配置中心**: 集中配置管理，支持动态更新
- **监控告警**: Prometheus + Grafana监控体系
- **链路追踪**: 分布式调用链追踪
- **负载均衡**: 多种负载均衡算法
- **服务治理**: 健康检查、故障转移

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **RPC框架**: gRPC
- **数据库**: PostgreSQL + Redis
- **服务发现**: etcd
- **监控**: Prometheus + Grafana
- **链路追踪**: Jaeger
- **容器化**: Docker + Kubernetes

## 📦 快速开始

### 1. 环境准备

```bash
# 安装依赖
go mod tidy

# 启动基础设施
docker-compose up -d
```

### 2. 启动服务

```bash
# 启动服务注册中心
go run infrastructure/service-registry/main.go

# 启动用户服务
go run services/user-service/main.go

# 启动API网关
go run api-gateway/main.go
```

### 3. 访问服务

- API网关: http://localhost:8080
- 监控面板: http://localhost:3000
- 服务注册中心: http://localhost:2379

## 📚 文档

详细文档请查看 [docs](./docs/) 目录。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License