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

            if [ -f "cmd/root.go" ]; then
                MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}' 2>/dev/null || echo "$PROJECT_NAME")
                CMD_NAME=$(basename "$MODULE_NAME" 2>/dev/null || echo "$PROJECT_NAME")

                cat > cmd/root.go << 'ROOT_EOF'

ROOT_EOF
                sed -i.bak "s|\$PROJECT_NAME|$MODULE_NAME|g" cmd/root.go
                sed -i.bak "s|\$CMD_NAME|$CMD_NAME|g" cmd/root.go
                rm -f cmd/root.go.bak
            fi
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
    git add . &> /dev/null
    git commit -m "chore: initial commit" &> /dev/null && echo "✓ 项目初始化提交"
fi

if [ -f ".pre-commit-config.yaml" ] && command -v pre-commit &> /dev/null; then
    pre-commit install &> /dev/null && echo "✓ 安装 pre-commit hooks"
fi

if [ -f "Makefile" ] && command -v make &> /dev/null && grep -q "^fmt:" Makefile; then
    make fmt &> /dev/null && echo "✓ 格式化代码"
fi
