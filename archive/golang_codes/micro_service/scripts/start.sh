#!/bin/bash

# 分布式微服务架构平台启动脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    log_info "Docker环境检查通过"
}

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.21+"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go版本: $GO_VERSION"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    mkdir -p bin
    mkdir -p logs
    mkdir -p data/postgres
    mkdir -p data/redis
    mkdir -p data/etcd
    mkdir -p data/prometheus
    mkdir -p data/grafana
}

# 安装Go依赖
install_dependencies() {
    log_info "安装Go依赖..."
    go mod download
    go mod tidy
    log_info "依赖安装完成"
}

# 构建服务
build_services() {
    log_info "构建服务..."
    
    # 构建API网关
    log_info "构建API网关..."
    cd api-gateway
    go build -o ../bin/api-gateway ./main.go
    cd ..
    
    # 构建用户服务
    log_info "构建用户服务..."
    cd services/user-service
    go build -o ../../bin/user-service ./main.go
    cd ../..
    
    log_info "服务构建完成"
}

# 启动基础设施服务
start_infrastructure() {
    log_info "启动基础设施服务..."
    docker-compose up -d postgres redis etcd prometheus grafana jaeger
    
    log_info "等待服务启动..."
    sleep 15
    
    # 检查服务状态
    check_service_health "PostgreSQL" "localhost" "5432"
    check_service_health "Redis" "localhost" "6379"
    check_service_health "etcd" "localhost" "2379"
}

# 检查服务健康状态
check_service_health() {
    local service_name=$1
    local host=$2
    local port=$3
    
    log_info "检查 $service_name 服务状态..."
    
    for i in {1..30}; do
        if nc -z $host $port 2>/dev/null; then
            log_info "$service_name 服务已就绪"
            return 0
        fi
        log_debug "等待 $service_name 服务启动... ($i/30)"
        sleep 2
    done
    
    log_error "$service_name 服务启动失败"
    return 1
}

# 运行数据库迁移
run_migrations() {
    log_info "运行数据库迁移..."
    
    # 等待PostgreSQL完全启动
    sleep 5
    
    # 执行初始化脚本
    if command -v psql &> /dev/null; then
        PGPASSWORD=password psql -h localhost -p 5432 -U postgres -f scripts/init-db.sql
        log_info "数据库迁移完成"
    else
        log_warn "psql未安装，跳过数据库迁移"
        log_warn "请手动执行: psql -h localhost -p 5432 -U postgres -f scripts/init-db.sql"
    fi
}

# 启动微服务
start_microservices() {
    log_info "启动微服务..."
    
    # 启动用户服务
    log_info "启动用户服务..."
    nohup ./bin/user-service > logs/user-service.log 2>&1 &
    USER_SERVICE_PID=$!
    echo $USER_SERVICE_PID > logs/user-service.pid
    
    sleep 3
    
    # 启动API网关
    log_info "启动API网关..."
    nohup ./bin/api-gateway > logs/api-gateway.log 2>&1 &
    API_GATEWAY_PID=$!
    echo $API_GATEWAY_PID > logs/api-gateway.pid
    
    sleep 3
    
    log_info "微服务启动完成"
}

# 检查微服务健康状态
check_microservices_health() {
    log_info "检查微服务健康状态..."
    
    # 检查用户服务
    if curl -f http://localhost:8081/health >/dev/null 2>&1; then
        log_info "用户服务健康检查通过"
    else
        log_error "用户服务健康检查失败"
    fi
    
    # 检查API网关
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        log_info "API网关健康检查通过"
    else
        log_error "API网关健康检查失败"
    fi
}

# 显示服务信息
show_service_info() {
    log_info "=== 服务信息 ==="
    echo -e "${GREEN}API网关:${NC} http://localhost:8080"
    echo -e "${GREEN}用户服务:${NC} http://localhost:8081"
    echo -e "${GREEN}Prometheus:${NC} http://localhost:9090"
    echo -e "${GREEN}Grafana:${NC} http://localhost:3000 (admin/admin)"
    echo -e "${GREEN}Jaeger:${NC} http://localhost:16686"
    echo -e "${GREEN}PostgreSQL:${NC} localhost:5432 (postgres/password)"
    echo -e "${GREEN}Redis:${NC} localhost:6379"
    echo -e "${GREEN}etcd:${NC} localhost:2379"
    echo ""
    log_info "=== 日志文件 ==="
    echo -e "${BLUE}API网关日志:${NC} logs/api-gateway.log"
    echo -e "${BLUE}用户服务日志:${NC} logs/user-service.log"
    echo ""
    log_info "=== 停止服务 ==="
    echo -e "${YELLOW}运行停止脚本:${NC} ./scripts/stop.sh"
}

# 主函数
main() {
    log_info "启动分布式微服务架构平台..."
    
    # 检查环境
    check_docker
    check_go
    
    # 创建目录
    create_directories
    
    # 安装依赖
    install_dependencies
    
    # 构建服务
    build_services
    
    # 启动基础设施
    start_infrastructure
    
    # 运行数据库迁移
    run_migrations
    
    # 启动微服务
    start_microservices
    
    # 等待服务完全启动
    sleep 10
    
    # 健康检查
    check_microservices_health
    
    # 显示服务信息
    show_service_info
    
    log_info "平台启动完成！"
}

# 错误处理
trap 'log_error "启动过程中发生错误，请检查日志"; exit 1' ERR

# 执行主函数
main "$@"