# 微服务状态检查脚本
Write-Host "检查微服务状态..." -ForegroundColor Green

# 检查进程状态
Write-Host "进程状态:" -ForegroundColor Yellow
$processes = Get-Process | Where-Object {$_.ProcessName -eq "user-service" -or $_.ProcessName -eq "api-gateway"}
if ($processes) {
    $processes | Select-Object ProcessName, Id, CPU | Format-Table
} else {
    Write-Host "没有找到微服务进程" -ForegroundColor Red
}

# 检查API网关健康状态
Write-Host "API网关健康检查:" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "API网关正常运行" -ForegroundColor Green
        Write-Host $response.Content -ForegroundColor Cyan
    }
} catch {
    Write-Host "API网关连接失败" -ForegroundColor Red
}

# 检查用户服务健康状态
Write-Host "用户服务健康检查:" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "用户服务正常运行" -ForegroundColor Green
        Write-Host $response.Content -ForegroundColor Cyan
    }
} catch {
    Write-Host "用户服务连接失败" -ForegroundColor Red
}

Write-Host "检查完成！" -ForegroundColor Green