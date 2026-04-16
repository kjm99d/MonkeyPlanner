@echo off
setlocal enabledelayedexpansion

:: MonkeyPlanner MCP auto-updater (Windows batch)

set "REPO=kjm99d/MonkeyPlanner"
set "BINARY_NAME=monkey-planner-windows-amd64.exe"
set "INSTALL_DIR=%~dp0"
set "BINARY_PATH=%INSTALL_DIR%%BINARY_NAME%"
set "VERSION_FILE=%INSTALL_DIR%.current-version"

:: Read current version
set "CURRENT_VERSION="
if exist "%VERSION_FILE%" (
  set /p CURRENT_VERSION=<"%VERSION_FILE%"
)

:: Check latest release via GitHub API
set "LATEST_TAG="
for /f "tokens=2 delims=:, " %%a in ('curl -sf "https://api.github.com/repos/%REPO%/releases/latest" 2^>nul ^| findstr "tag_name"') do (
  set "LATEST_TAG=%%~a"
)

if not defined LATEST_TAG (
  echo [MCP] Could not check for updates, using current version >&2
  goto :run
)

if "%LATEST_TAG%"=="%CURRENT_VERSION%" (
  echo [MCP] Already up to date: %CURRENT_VERSION% >&2
  goto :run
)

echo [MCP] Updating %CURRENT_VERSION% -^> %LATEST_TAG% >&2

set "DOWNLOAD_URL=https://github.com/%REPO%/releases/download/%LATEST_TAG%/%BINARY_NAME%"
set "TEMP_PATH=%INSTALL_DIR%%BINARY_NAME%.tmp"

curl -sfL "%DOWNLOAD_URL%" -o "%TEMP_PATH%" 2>nul
if %errorlevel% neq 0 (
  echo [MCP] Download failed, using current version >&2
  del /f "%TEMP_PATH%" 2>nul
  goto :run
)

:: Backup and replace
if exist "%BINARY_PATH%" (
  move /y "%BINARY_PATH%" "%INSTALL_DIR%%BINARY_NAME%.bak" >nul 2>&1
)
move /y "%TEMP_PATH%" "%BINARY_PATH%" >nul 2>&1
echo %LATEST_TAG%>"%VERSION_FILE%"
echo [MCP] Updated to %LATEST_TAG% >&2

:run
"%BINARY_PATH%" %*
