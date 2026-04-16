#!/bin/bash
# MonkeyPlanner MCP auto-updater
# GitHub 최신 릴리즈를 체크하고, 새 버전이면 다운로드 후 실행

set -euo pipefail

REPO="kjm99d/MonkeyPlanner"
BINARY_NAME="monkey-planner-windows-amd64.exe"
INSTALL_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"
VERSION_FILE="$INSTALL_DIR/.current-version"

# 현재 버전 읽기
CURRENT_VERSION=""
if [ -f "$VERSION_FILE" ]; then
  CURRENT_VERSION=$(cat "$VERSION_FILE")
fi

# GitHub 최신 릴리즈 체크 (stderr로 로그, stdout은 MCP 프로토콜 전용)
LATEST=$(curl -sf "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null || echo "")

if [ -n "$LATEST" ]; then
  LATEST_TAG=$(echo "$LATEST" | grep -o '"tag_name": *"[^"]*"' | head -1 | cut -d'"' -f4)

  if [ -n "$LATEST_TAG" ] && [ "$LATEST_TAG" != "$CURRENT_VERSION" ]; then
    echo "Updating $CURRENT_VERSION -> $LATEST_TAG" >&2

    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$BINARY_NAME"
    TEMP_PATH="$INSTALL_DIR/${BINARY_NAME}.tmp"

    if curl -sfL "$DOWNLOAD_URL" -o "$TEMP_PATH" 2>/dev/null; then
      # 기존 바이너리 백업
      if [ -f "$BINARY_PATH" ]; then
        mv "$BINARY_PATH" "$INSTALL_DIR/${BINARY_NAME}.bak"
      fi
      mv "$TEMP_PATH" "$BINARY_PATH"
      echo "$LATEST_TAG" > "$VERSION_FILE"
      echo "Updated to $LATEST_TAG" >&2
    else
      echo "Download failed, using current version" >&2
      rm -f "$TEMP_PATH"
    fi
  else
    echo "Already up to date: $CURRENT_VERSION" >&2
  fi
else
  echo "Could not check for updates, using current version" >&2
fi

# MCP 서버 실행
exec "$BINARY_PATH" "$@"
