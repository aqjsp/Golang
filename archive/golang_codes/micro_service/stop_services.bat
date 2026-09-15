@echo off
echo 🛑 停止分布式微服务架构平台
echo.

echo 🔍 查找并停止Go服务进程...

:: 停止用户服务
for /f "tokens=2" %%i in ('tasklist /fi "imagename eq user-service.exe" /fo csv ^| find "user-service.exe"') do (
    echo 停止用户服务进程 %%i
    taskkill /pid %%i /f
)

:: 停止API网关
for /f "tokens=2" %%i in ('tasklist /fi "imagename eq api-gateway.exe" /fo csv ^| find "api-gateway.exe"') do (
    echo 停止API网关进程 %%i
    taskkill /pid %%i /f
)

:: 停止可能的Go进程
for /f "tokens=2" %%i in ('tasklist /fi "imagename eq go.exe" /fo csv ^| find "go.exe"') do (
    echo 停止Go进程 %%i
    taskkill /pid %%i /f
)

echo.
echo ✅ 服务停止完成！
echo.
echo 📝 清理信息：
echo    - 所有微服务进程已停止
echo    - 日志文件保留在 logs/ 目录
echo    - 可以使用 start_without_docker.bat 重新启动
echo.
pause