#!/bin/bash

# 构建脚本 - 自动注入版本信息

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取版本信息
get_version_info() {
    # 优先从VERSION文件读取，然后尝试git tag，最后使用默认值
    if [ -f "VERSION" ]; then
        VERSION=$(cat VERSION | tr -d '\n')
    else
        VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.8.0")
    fi
    
    # 获取git commit hash (短版本)
    GIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    # 获取git状态（是否有未提交的更改）
    if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
        GIT_HASH="${GIT_HASH}-dirty"
    fi
    
    # 获取构建时间
    BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
    
    echo -e "${GREEN}版本信息:${NC}"
    echo "  Version:    $VERSION"
    echo "  Git Hash:   $GIT_HASH"
    echo "  Build Time: $BUILD_TIME"
}

# 构建单个平台
build_single() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3
    
    echo -e "${YELLOW}正在构建 ${GOOS}/${GOARCH}...${NC}"
    
    # 构建命令，注入版本信息
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -trimpath \
        -ldflags="-w -s \
            -X 'main.version=${VERSION}' \
            -X 'main.gitHash=${GIT_HASH}' \
            -X 'main.buildTime=${BUILD_TIME}'" \
        -o "$OUTPUT" \
        ./src/main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 构建成功: $OUTPUT${NC}"
        return 0
    else
        echo -e "${RED}✗ 构建失败: ${GOOS}/${GOARCH}${NC}"
        return 1
    fi
}

# 主函数
main() {
    echo -e "${GREEN}=== TCPing 构建脚本 ===${NC}\n"
    
    # 获取版本信息
    get_version_info
    echo ""
    
    # 检查是否在正确的目录
    if [ ! -f "src/main.go" ]; then
        echo -e "${RED}错误: 未找到 src/main.go${NC}"
        echo "请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 默认构建当前平台
    if [ $# -eq 0 ]; then
        echo -e "${YELLOW}构建当前平台...${NC}"
        OUTPUT="tcping"
        
        # 自动检测当前平台
        CURRENT_OS=$(go env GOOS)
        CURRENT_ARCH=$(go env GOARCH)
        
        # Windows需要.exe后缀
        if [ "$CURRENT_OS" = "windows" ]; then
            OUTPUT="tcping.exe"
        fi
        
        build_single "$CURRENT_OS" "$CURRENT_ARCH" "$OUTPUT"
        
        if [ $? -eq 0 ]; then
            # 显示文件信息
            ls -lh "$OUTPUT"
            
            # 测试版本输出
            echo -e "\n${GREEN}版本信息测试:${NC}"
            ./"$OUTPUT" -V
        fi
    else
        # 构建指定平台
        case "$1" in
            all)
                # 构建所有平台（调用compiler.sh）
                echo -e "${YELLOW}构建所有平台...${NC}"
                echo "调用 compiler.sh 进行跨平台批量编译..."
                chmod +x compiler.sh
                ./compiler.sh
                ;;
            linux)
                build_single linux amd64 "tcping-linux-amd64"
                ;;
            windows)
                build_single windows amd64 "tcping-windows-amd64.exe"
                ;;
            darwin|macos)
                build_single darwin amd64 "tcping-darwin-amd64"
                ;;
            test)
                # 仅测试版本信息
                echo -e "${YELLOW}测试版本信息...${NC}"
                get_version_info
                ;;
            *)
                echo "用法: $0 [all|linux|windows|darwin|test]"
                echo ""
                echo "选项:"
                echo "  不带参数: 构建当前平台的二进制文件"
                echo "  all:      构建所有平台（调用compiler.sh）"  
                echo "  linux:    构建Linux amd64"
                echo "  windows:  构建Windows amd64"
                echo "  darwin:   构建macOS amd64"
                echo "  test:     仅显示版本信息"
                echo ""
                echo "说明:"
                echo "  - 单平台构建用于快速测试"
                echo "  - 跨平台批量编译请使用 './build.sh all' 或直接运行 './compiler.sh'"
                exit 1
                ;;
        esac
    fi
}

# 运行主函数
main "$@"