@echo off
setlocal enabledelayedexpansion

:: 创建输出目录
if not exist dist (
    mkdir dist
)

:: 可执行文件名称
set app_name=gocs_bot

:: 构建 Windows 下的可执行文件
echo Building for Windows...
go build -o dist\%app_name%.exe main.go

if errorlevel 1 (
    echo An error has occurred! Aborting the script execution...
    exit /b 1
)

:: 构建 Linux 下的可执行文件
echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o dist\%app_name% main.go


echo Build complete. Executables are located in the 'dist' directory.
