@echo off
setlocal

:: Путь к папке с бэкендом
set "backend_dir=%~dp0back"
:: Путь к папке с фронтендом
set "frontend_dir=%~dp0front"

:: Запуск фронтенда в новом окне
cd /d "%frontend_dir%"
start cmd /k "npm run dev"

:: Инициализация Swagger и запуск бэкенда в одном окне
cd /d "%backend_dir%"
swag init && cmd /k "go run main.go"



pause
endlocal