@echo off
echo 🚀 启动分布式微服务架构平台（无Docker模式）
echo.

echo 📝 注意：由于未安装Docker，将以简化模式启动服务
echo    - 数据库连接将失败（正常现象）
echo    - 服务注册将失败（正常现象）
echo    - 但可以验证服务基本功能
echo.

echo 🔧 检查Go环境...
go version
if %errorlevel% neq 0 (
    echo ❌ Go环境未正确配置
    pause
    exit /b 1
)

echo ✅ Go环境正常
echo.

echo 🏗️ 构建服务...
go build -o bin/user-service.exe ./services/user-service/main.go
if %errorlevel% neq 0 (
    echo ❌ 用户服务构建失败
    pause
    exit /b 1
)

go build -o bin/api-gateway.exe ./api-gateway/main.go
if %errorlevel% neq 0 (
    echo ❌ API网关构建失败
    pause
    exit /b 1
)

echo ✅ 服务构建完成
echo.

echo 🚀 启动用户服务（后台运行）...
start "用户服务" /min cmd /c "bin\user-service.exe > logs\user-service.log 2>&1"

echo 等待用户服务启动...
timeout /t 3 /nobreak > nul

echo 🚀 启动API网关（后台运行）...
start "API网关" /min cmd /c "bin\api-gateway.exe > logs\api-gateway.log 2>&1"

echo 等待API网关启动...
timeout /t 3 /nobreak > nul

echo.
echo ✅ 服务启动完成！
echo.
echo 📊 运行健康检查测试...
go run test_project.go

echo.
echo 🎯 启动完成！
echo.
echo 📝 服务信息：
echo    - API网关: http://localhost:8080
echo    - 用户服务: http://localhost:8081
echo    - 日志文件: logs/ 目录
echo.
echo 💡 提示：
echo    - 使用 stop_services.bat 停止所有服务
echo    - 查看 logs/ 目录中的日志文件了解详细信息
echo    - 安装Docker后可使用完整功能
echo.
pause