#!/bin/bash

# 分布式微服务架构平台停止脚本

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

# 停止微服务进程
stop_microservices() {
    log_info "停止微服务..."
    
    # 停止API网关
    if [ -f "logs/api-gateway.pid" ]; then
        API_GATEWAY_PID=$(cat logs/api-gateway.pid)
        if ps -p $API_GATEWAY_PID > /dev/null 2>&1; then
            log_info "停止API网关 (PID: $API_GATEWAY_PID)..."
            kill $API_GATEWAY_PID
            sleep 2
            
            # 强制杀死进程（如果还在运行）
            if ps -p $API_GATEWAY_PID > /dev/null 2>&1; then
                log_warn "强制停止API网关..."
                kill -9 $API_GATEWAY_PID
            fi
            
            rm -f logs/api-gateway.pid
            log_info "API网关已停止"
        else
            log_warn "API网关进程不存在"
            rm -f logs/api-gateway.pid
        fi
    else
        log_warn "未找到API网关PID文件"
    fi
    
    # 停止用户服务
    if [ -f "logs/user-service.pid" ]; then
        USER_SERVICE_PID=$(cat logs/user-service.pid)
        if ps -p $USER_SERVICE_PID > /dev/null 2>&1; then
            log_info "停止用户服务 (PID: $USER_SERVICE_PID)..."
            kill $USER_SERVICE_PID
            sleep 2
            
            # 强制杀死进程（如果还在运行）
            if ps -p $USER_SERVICE_PID > /dev/null 2>&1; then
                log_warn "强制停止用户服务..."
                kill -9 $USER_SERVICE_PID
            fi
            
            rm -f logs/user-service.pid
            log_info "用户服务已停止"
        else
            log_warn "用户服务进程不存在"
            rm -f logs/user-service.pid
        fi
    else
        log_warn "未找到用户服务PID文件"
    fi
    
    # 停止订单服务（如果存在）
    if [ -f "logs/order-service.pid" ]; then
        ORDER_SERVICE_PID=$(cat logs/order-service.pid)
        if ps -p $ORDER_SERVICE_PID > /dev/null 2>&1; then
            log_info "停止订单服务 (PID: $ORDER_SERVICE_PID)..."
            kill $ORDER_SERVICE_PID
            sleep 2
            
            if ps -p $ORDER_SERVICE_PID > /dev/null 2>&1; then
                log_warn "强制停止订单服务..."
                kill -9 $ORDER_SERVICE_PID
            fi
            
            rm -f logs/order-service.pid
            log_info "订单服务已停止"
        else
            log_warn "订单服务进程不存在"
            rm -f logs/order-service.pid
        fi
    fi
    
    # 停止支付服务（如果存在）
    if [ -f "logs/payment-service.pid" ]; then
        PAYMENT_SERVICE_PID=$(cat logs/payment-service.pid)
        if ps -p $PAYMENT_SERVICE_PID > /dev/null 2>&1; then
            log_info "停止支付服务 (PID: $PAYMENT_SERVICE_PID)..."
            kill $PAYMENT_SERVICE_PID
            sleep 2
            
            if ps -p $PAYMENT_SERVICE_PID > /dev/null 2>&1; then
                log_warn "强制停止支付服务..."
                kill -9 $PAYMENT_SERVICE_PID
            fi
            
            rm -f logs/payment-service.pid
            log_info "支付服务已停止"
        else
            log_warn "支付服务进程不存在"
            rm -f logs/payment-service.pid
        fi
    fi
    
    log_info "微服务停止完成"
}

# 停止Docker容器
stop_docker_services() {
    log_info "停止Docker服务..."
    
    if command -v docker-compose &> /dev/null; then
        if [ -f "docker-compose.yml" ]; then
            docker-compose down
            log_info "Docker服务已停止"
        else
            log_warn "未找到docker-compose.yml文件"
        fi
    else
        log_warn "Docker Compose未安装，跳过Docker服务停止"
    fi
}

# 清理进程（通过端口查找）
cleanup_by_port() {
    log_info "清理端口占用的进程..."
    
    local ports=("8080" "8081" "8082" "8083")
    
    for port in "${ports[@]}"; do
        local pid=$(lsof -ti:$port 2>/dev/null || true)
        if [ ! -z "$pid" ]; then
            log_info "发现端口 $port 被进程 $pid 占用，正在停止..."
            kill $pid 2>/dev/null || true
            sleep 1
            
            # 检查进程是否还在运行
            if ps -p $pid > /dev/null 2>&1; then
                log_warn "强制停止进程 $pid..."
                kill -9 $pid 2>/dev/null || true
            fi
            
            log_info "端口 $port 已释放"
        fi
    done
}

# 清理临时文件
cleanup_temp_files() {
    log_info "清理临时文件..."
    
    # 清理PID文件
    rm -f logs/*.pid
    
    # 清理临时日志（可选）
    # rm -f logs/*.log
    
    log_info "临时文件清理完成"
}

# 显示停止后的状态
show_status() {
    log_info "=== 停止状态 ==="
    
    # 检查端口是否还被占用
    local ports=("8080" "8081" "8082" "8083")
    local occupied_ports=()
    
    for port in "${ports[@]}"; do
        if lsof -ti:$port >/dev/null 2>&1; then
            occupied_ports+=("$port")
        fi
    done
    
    if [ ${#occupied_ports[@]} -eq 0 ]; then
        log_info "所有服务端口已释放"
    else
        log_warn "以下端口仍被占用: ${occupied_ports[*]}"
        log_warn "请手动检查并停止相关进程"
    fi
    
    # 检查Docker容器状态
    if command -v docker &> /dev/null; then
        local running_containers=$(docker ps --filter "name=microservice_" --format "table {{.Names}}\t{{.Status}}" 2>/dev/null || true)
        if [ ! -z "$running_containers" ] && [ "$running_containers" != "NAMES	STATUS" ]; then
            log_warn "以下Docker容器仍在运行:"
            echo "$running_containers"
        else
            log_info "所有相关Docker容器已停止"
        fi
    fi
}

# 强制清理所有相关进程
force_cleanup() {
    log_warn "执行强制清理..."
    
    # 通过进程名查找并杀死相关进程
    local process_names=("api-gateway" "user-service" "order-service" "payment-service")
    
    for process_name in "${process_names[@]}"; do
        local pids=$(pgrep -f $process_name 2>/dev/null || true)
        if [ ! -z "$pids" ]; then
            log_info "强制停止 $process_name 进程: $pids"
            echo $pids | xargs kill -9 2>/dev/null || true
        fi
    done
    
    # 清理端口占用
    cleanup_by_port
    
    log_info "强制清理完成"
}

# 主函数
main() {
    log_info "停止分布式微服务架构平台..."
    
    # 解析命令行参数
    local force_mode=false
    local keep_docker=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--force)
                force_mode=true
                shift
                ;;
            --keep-docker)
                keep_docker=true
                shift
                ;;
            -h|--help)
                echo "用法: $0 [选项]"
                echo "选项:"
                echo "  -f, --force       强制停止所有相关进程"
                echo "  --keep-docker     保持Docker服务运行"
                echo "  -h, --help        显示帮助信息"
                exit 0
                ;;
            *)
                log_error "未知选项: $1"
                exit 1
                ;;
        esac
    done
    
    # 停止微服务
    stop_microservices
    
    # 停止Docker服务（除非指定保持运行）
    if [ "$keep_docker" = false ]; then
        stop_docker_services
    else
        log_info "保持Docker服务运行"
    fi
    
    # 清理端口占用
    cleanup_by_port
    
    # 强制清理（如果指定）
    if [ "$force_mode" = true ]; then
        force_cleanup
    fi
    
    # 清理临时文件
    cleanup_temp_files
    
    # 显示状态
    show_status
    
    log_info "平台停止完成！"
}

# 错误处理
trap 'log_error "停止过程中发生错误"; exit 1' ERR

# 执行主函数
main "$@"