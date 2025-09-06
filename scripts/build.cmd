@echo off
REM Windows batch build script for skv binary
setlocal enabledelayedexpansion

set MODE=%1
if "%MODE%"=="" set MODE=host

if "%MODE%"=="host" goto build_host
if "%MODE%"=="all" goto build_all
echo Usage: %0 [host^|all]
echo   host (default): build for current platform
echo   all: build for multiple platforms
exit /b 2

:build_host
echo Building for current platform...
for /f %%i in ('go env GOOS') do set GOOS=%%i
for /f %%i in ('go env GOARCH') do set GOARCH=%%i
call :build_target !GOOS! !GOARCH!
goto end

:build_all
echo Building for multiple platforms...
call :build_target darwin arm64
call :build_target darwin amd64
call :build_target linux amd64
call :build_target windows amd64
goto end

:build_target
set TARGET_GOOS=%1
set TARGET_GOARCH=%2
echo Building skv for %TARGET_GOOS%/%TARGET_GOARCH%...

if not exist dist mkdir dist

REM Get version info
for /f %%i in ('git describe --tags --dirty --always 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

for /f %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT=%%i
if "%COMMIT%"=="" set COMMIT=

REM Get current timestamp in ISO format
for /f "tokens=1-3 delims=/ " %%a in ('date /t') do set DATE_PART=%%c-%%a-%%b
for /f "tokens=1-2 delims=: " %%a in ('time /t') do set TIME_PART=%%a:%%b:00
set BUILD_DATE=%DATE_PART%T%TIME_PART%Z

set OUTPUT_NAME=skv_%TARGET_GOOS%_%TARGET_GOARCH%
if "%TARGET_GOOS%"=="windows" set OUTPUT_NAME=%OUTPUT_NAME%.exe

set LDFLAGS=-s -w -X skv/internal/version.Version=%VERSION% -X skv/internal/version.Commit=%COMMIT% -X skv/internal/version.Date=%BUILD_DATE%

set GOOS=%TARGET_GOOS%
set GOARCH=%TARGET_GOARCH%
go build -trimpath -ldflags "%LDFLAGS%" -o dist\%OUTPUT_NAME% ./cmd/skv

goto :eof

:end
echo Build completed. Artifacts in dist/
dir /b dist\
endlocal
