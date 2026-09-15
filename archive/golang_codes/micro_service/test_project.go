package main

import (
	"fmt"
	"net/http"
	"time"
)

// 简单的HTTP健康检查测试
func testHTTPEndpoint(url string, name string) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("❌ %s: 连接失败 - %v\n", name, err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("✅ %s: 服务正常运行 (状态码: %d)\n", name, resp.StatusCode)
	} else {
		fmt.Printf("⚠️  %s: 服务响应异常 (状态码: %d)\n", name, resp.StatusCode)
	}
}

func main() {
	fmt.Println("🚀 开始测试微服务架构平台...")
	fmt.Println("")
	
	// 测试API网关
	fmt.Println("📡 测试API网关:")
	testHTTPEndpoint("http://localhost:8080/health", "API网关健康检查")
	testHTTPEndpoint("http://localhost:8080/api/v1/services", "服务发现")
	
	fmt.Println("")
	
	// 测试用户服务
	fmt.Println("👤 测试用户服务:")
	testHTTPEndpoint("http://localhost:8081/health", "用户服务健康检查")
	testHTTPEndpoint("http://localhost:8081/metrics", "用户服务指标")
	
	fmt.Println("")
	
	// 测试监控服务
	fmt.Println("📊 测试监控服务:")
	testHTTPEndpoint("http://localhost:9090", "Prometheus")
	testHTTPEndpoint("http://localhost:3000", "Grafana")
	testHTTPEndpoint("http://localhost:16686", "Jaeger")
	
	fmt.Println("")
	fmt.Println("✨ 测试完成！")
	fmt.Println("")
	fmt.Println("📝 注意事项:")
	fmt.Println("   - 如果服务显示连接失败，请确保相应的服务已启动")
	fmt.Println("   - 完整的功能测试需要启动所有基础设施服务（PostgreSQL, Redis, etcd等）")
	fmt.Println("   - 可以使用 docker-compose up -d 启动所有依赖服务")
}