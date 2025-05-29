#!/bin/bash

# 源代码路径
SRC_PATH="./src/main.go"
# 输出目录
OUT_DIR="./bin"
# 程序名
APP_NAME="tcping"

# 需要编译的目标平台和架构
PLATFORMS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "linux/loong64"
  "windows/386"
  "windows/amd64"
  "windows/arm"
  "windows/arm64"
)

# 清理之前的编译产物
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

# 初始化 SHA256SUMS 文件
> "$OUT_DIR/SHA256SUMS.txt"

# 编译每个平台并计算 SHA256
for PLATFORM in "${PLATFORMS[@]}"; do
  # 获取平台的 GOOS 和 GOARCH
  GOOS=$(echo "$PLATFORM" | cut -d'/' -f1)
  GOARCH=$(echo "$PLATFORM" | cut -d'/' -f2)

  # 设置输出文件路径，确保文件名包含平台信息
  OUT_FILE="$OUT_DIR/${APP_NAME}-${GOOS}-${GOARCH}"

  # 设置环境变量并编译
  echo "编译 ${GOOS}/${GOARCH}..."
  CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags="-w -s" -o "$OUT_FILE" "$SRC_PATH"

  # 判断是否是 Windows 平台，需要添加 .exe 扩展名
  if [ "$GOOS" == "windows" ]; then
    mv "$OUT_FILE" "$OUT_FILE.exe"
    OUT_FILE="$OUT_FILE.exe"
  fi

  # 计算 SHA256 值并追加到 SHA256SUMS.txt 中，同时输出到终端
  echo "计算 SHA256..."
  sha256sum "$OUT_FILE" | tee -a "$OUT_DIR/SHA256SUMS.txt"

done

echo "编译完成，所有文件已存储在 $OUT_DIR 目录下。"
