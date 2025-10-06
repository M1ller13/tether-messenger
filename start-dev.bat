@echo off
echo Starting Tether Messenger Development Environment...

echo.
echo 1. Starting PostgreSQL (if not running)...
echo Please make sure PostgreSQL is installed and running on localhost:5432
echo Database: tether_messenger
echo User: postgres
echo Password: password

echo.
echo 2. Starting Backend Server...
cd server
start "Tether Backend" cmd /k "go run main.go"

echo.
echo 3. Starting Frontend Server...
cd ..\client
start "Tether Frontend" cmd /k "npm run dev"

echo.
echo Development environment started!
echo.
echo Backend: http://localhost:8081
echo Frontend: http://localhost:3000
echo.
echo Press any key to exit...
pause > nul
