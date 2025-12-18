#!/bin/bash

PROJECT_NAME="${GOGENIE_PROJECT_NAME:-$(basename "$(pwd)")}"

if [ -f "go.mod" ]; then
    go mod tidy &> /dev/null && echo "✓ 更新依赖"
    go mod download &> /dev/null
else
    if go mod init "$PROJECT_NAME" &> /dev/null; then
        echo "✓ 初始化 Go 模块"

        AUTHOR=$(git config user.name 2>/dev/null || echo "Unknown")
        if cobra-cli init --author "$AUTHOR" --license apache &> /dev/null || cobra-cli init &> /dev/null; then
            echo "✓ 初始化 Cobra CLI"
        fi

        cobra-cli add server &> /dev/null && echo "✓ 添加 server 命令"
        go mod tidy &> /dev/null
    fi
fi

if command -v gogenie &> /dev/null; then
    gogenie init -d &> /dev/null && echo "✓ 初始化项目目录"
fi

if [ ! -d ".git" ]; then
    git init &> /dev/null && echo "✓ 初始化 Git"
    git checkout -b main &> /dev/null || git checkout -b master &> /dev/null || true
fi

if [ -f ".pre-commit-config.yaml" ] && command -v pre-commit &> /dev/null; then
    pre-commit install &> /dev/null && echo "✓ 安装 pre-commit hooks"
fi

if [ -f "Makefile" ] && command -v make &> /dev/null && grep -q "^fmt:" Makefile; then
    make fmt &> /dev/null && echo "✓ 格式化代码"
fi
