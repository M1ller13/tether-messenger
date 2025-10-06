Write-Host "Starting Tether Messenger Development Environment..." -ForegroundColor Green

Write-Host ""
Write-Host "1. Starting PostgreSQL (if not running)..." -ForegroundColor Yellow
Write-Host "Please make sure PostgreSQL is installed and running on localhost:5432" -ForegroundColor Cyan
Write-Host "Database: tether_messenger" -ForegroundColor Cyan
Write-Host "User: postgres" -ForegroundColor Cyan
Write-Host "Password: password" -ForegroundColor Cyan

Write-Host ""
Write-Host "2. Starting Backend Server..." -ForegroundColor Yellow
Set-Location server
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run main.go"

Write-Host ""
Write-Host "3. Starting Frontend Server..." -ForegroundColor Yellow
Set-Location ..\client
Start-Process powershell -ArgumentList "-NoExit", "-Command", "npm run dev"

Write-Host ""
Write-Host "Development environment started!" -ForegroundColor Green
Write-Host ""
Write-Host "Backend: http://localhost:8081" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press any key to exit..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
