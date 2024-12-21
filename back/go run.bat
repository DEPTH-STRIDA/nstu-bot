@echo off
setlocal
set "batch_dir=%~dp0"
cd %batch_dir%
swag init
go run main.go
pause
endlocal
