# 分布式微服务架构平台 Makefile

# 变量定义
GO_VERSION := 1.21
DOCKER_REGISTRY := localhost:5000
PROJECT_NAME := micro-service-platform
VERSION := v1.0.0

# 服务列表
SERVICES := api-gateway user-service order-service payment-service

# 默认目标
.PHONY: all
all: build

# 帮助信息
.PHONY: help
help:
	@echo "分布式微服务架构平台构建脚本"
	@echo ""
	@echo "可用命令:"
	@echo "  help          显示帮助信息"
	@echo "  deps          安装依赖"
	@echo "  build         构建所有服务"
	@echo "  build-service 构建指定服务 (make build-service SERVICE=user-service)"
	@echo "  test          运行测试"
	@echo "  test-coverage 运行测试并生成覆盖率报告"
	@echo "  lint          代码检查"
	@echo "  fmt           格式化代码"
	@echo "  clean         清理构建文件"
	@echo "  docker-build  构建Docker镜像"
	@echo "  docker-push   推送Docker镜像"
	@echo "  docker-up     启动Docker Compose"
	@echo "  docker-down   停止Docker Compose"
	@echo "  dev           启动开发环境"
	@echo "  migrate       运行数据库迁移"
	@echo ""

# 安装依赖
.PHONY: deps
deps:
	@echo "安装Go依赖..."
	go mod download
	go mod tidy
	@echo "依赖安装完成"

# 构建所有服务
.PHONY: build
build: deps
	@echo "构建所有服务..."
	@for service in $(SERVICES); do \
		echo "构建 $$service..."; \
		if [ "$$service" = "api-gateway" ]; then \
			cd api-gateway && go build -o ../bin/api-gateway ./main.go && cd ..; \
		else \
			cd services/$$service && go build -o ../../bin/$$service ./main.go && cd ../..; \
		fi; \
	done
	@echo "所有服务构建完成"

# 构建指定服务
.PHONY: build-service
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "请指定服务名称: make build-service SERVICE=user-service"; \
		exit 1; \
	fi
	@echo "构建服务: $(SERVICE)"
	@if [ "$(SERVICE)" = "api-gateway" ]; then \
		cd api-gateway && go build -o ../bin/api-gateway ./main.go; \
	else \
		cd services/$(SERVICE) && go build -o ../../bin/$(SERVICE) ./main.go; \
	fi
	@echo "服务 $(SERVICE) 构建完成"

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	go test -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 代码检查
.PHONY: lint
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，请先安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	go fmt ./...
	goimports -w .

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f
	@echo "清理完成"

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	@for service in $(SERVICES); do \
		echo "构建 $$service 镜像..."; \
		if [ "$$service" = "api-gateway" ]; then \
			docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f ./api-gateway/Dockerfile .; \
		else \
			docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f ./services/$$service/Dockerfile .; \
		fi; \
	done
	@echo "Docker镜像构建完成"

# 推送Docker镜像
.PHONY: docker-push
docker-push: docker-build
	@echo "推送Docker镜像..."
	@for service in $(SERVICES); do \
		echo "推送 $$service 镜像..."; \
		docker push $(DOCKER_REGISTRY)/$$service:$(VERSION); \
	done
	@echo "Docker镜像推送完成"

# 启动Docker Compose
.PHONY: docker-up
docker-up:
	@echo "启动Docker Compose..."
	docker-compose up -d
	@echo "服务启动完成"
	@echo "API网关: http://localhost:8080"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Jaeger: http://localhost:16686"

# 停止Docker Compose
.PHONY: docker-down
docker-down:
	@echo "停止Docker Compose..."
	docker-compose down
	@echo "服务已停止"

# 启动开发环境
.PHONY: dev
dev:
	@echo "启动开发环境..."
	@echo "启动基础设施服务..."
	docker-compose up -d postgres redis etcd prometheus grafana jaeger
	@echo "等待服务启动..."
	sleep 10
	@echo "开发环境已启动"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "etcd: localhost:2379"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000"
	@echo "Jaeger: http://localhost:16686"

# 运行数据库迁移
.PHONY: migrate
migrate:
	@echo "运行数据库迁移..."
	@if command -v psql >/dev/null 2>&1; then \
		psql -h localhost -p 5432 -U postgres -f scripts/init-db.sql; \
	else \
		echo "psql 未安装，请手动执行数据库迁移脚本"; \
	fi

# 生成API文档
.PHONY: docs
docs:
	@echo "生成API文档..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g api-gateway/main.go -o docs/; \
	else \
		echo "swag 未安装，请先安装: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 运行性能测试
.PHONY: benchmark
benchmark:
	@echo "运行性能测试..."
	go test -bench=. -benchmem ./...

# 检查服务健康状态
.PHONY: health-check
health-check:
	@echo "检查服务健康状态..."
	@curl -f http://localhost:8080/health || echo "API网关不可用"
	@curl -f http://localhost:8081/health || echo "用户服务不可用"
	@curl -f http://localhost:8082/health || echo "订单服务不可用"
	@curl -f http://localhost:8083/health || echo "支付服务不可用"

# 查看日志
.PHONY: logs
logs:
	@echo "查看服务日志..."
	docker-compose logs -f

# 重启服务
.PHONY: restart
restart: docker-down docker-up

# 安装开发工具
.PHONY: install-tools
install-tools:
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "开发工具安装完成"