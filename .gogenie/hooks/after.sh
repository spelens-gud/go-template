#!/bin/bash

echo ""
echo "=========================================="
echo "项目创建完成，开始初始化..."
echo "=========================================="
echo ""

PROJECT_NAME="${GOGENIE_PROJECT_NAME:-$(basename "$(pwd)")}"
PROJECT_PATH="${GOGENIE_PROJECT_PATH:-$(pwd)}"

echo "项目名称: $PROJECT_NAME"
echo "项目路径: $PROJECT_PATH"
echo ""

if [ -f "go.mod" ]; then
    echo "✓ 检测到 go.mod 文件"

    echo "正在更新 Go 模块依赖..."
    if go mod tidy; then
        echo "✓ Go 模块依赖更新完成"
    else
        echo "✗ Go 模块依赖更新失败"
    fi

    echo "正在下载 Go 模块依赖..."
    if go mod download; then
        echo "✓ Go 模块依赖下载完成"
    else
        echo "✗ Go 模块依赖下载失败"
    fi
else
    echo "未检测到 go.mod 文件，开始初始化 Go 模块和 CLI 项目..."

    echo "正在初始化 Go 模块: $PROJECT_NAME"
    if go mod init "$PROJECT_NAME"; then
        echo "✓ Go 模块初始化完成"
    else
        echo "✗ Go 模块初始化失败"
        exit 1
    fi

    echo "正在初始化 Cobra CLI 项目..."
    AUTHOR=$(git config user.name 2>/dev/null || echo "Unknown")
    if cobra-cli init --author "$AUTHOR" --license apache; then
        echo "✓ Cobra CLI 项目初始化完成"
    else
        echo "⚠ Cobra CLI 项目初始化失败，尝试使用默认参数..."
        if cobra-cli init; then
            echo "✓ Cobra CLI 项目初始化完成（使用默认参数）"
        else
            echo "✗ Cobra CLI 项目初始化失败，请检查错误信息"
            exit 1
        fi
    fi

    echo "正在添加 server 子命令..."
    if cobra-cli add server; then
        echo "✓ server 子命令添加完成"
    else
        echo "✗ server 子命令添加失败，请检查错误信息"
        exit 1
    fi

    echo "正在更新 Go 模块依赖..."
    if go mod tidy; then
        echo "✓ Go 模块依赖更新完成"
    else
        echo "✗ Go 模块依赖更新失败"
    fi
fi

echo ""

if [ ! -d ".git" ]; then
    echo "正在初始化 Git 仓库..."
    if git init; then
        echo "✓ Git 仓库初始化完成"

        if git config --global init.defaultBranch &> /dev/null; then
            git checkout -b main 2>/dev/null || git checkout -b master 2>/dev/null || true
        fi
    else
        echo "✗ Git 仓库初始化失败"
    fi
else
    echo "✓ Git 仓库已存在，跳过初始化"
fi

echo ""

if [ -f ".pre-commit-config.yaml" ]; then
    echo "检测到 .pre-commit-config.yaml 文件"

    if command -v pre-commit &> /dev/null; then
        echo "正在安装 pre-commit hooks..."
        if pre-commit install; then
            echo "✓ pre-commit hooks 安装完成"

            echo "正在运行 pre-commit 检查..."
            if pre-commit run --all-files; then
                echo "✓ pre-commit 检查完成"
            else
                echo "⚠ pre-commit 检查发现问题，请查看上方输出"
            fi
        else
            echo "✗ pre-commit hooks 安装失败"
        fi
    else
        echo "⚠ 未找到 pre-commit 命令，跳过 hooks 安装"
        echo "  提示: 请先安装 pre-commit: pip3 install pre-commit 或 brew install pre-commit"
    fi
else
    echo "提示: 未检测到 .pre-commit-config.yaml 文件，跳过 pre-commit 安装"
fi

echo ""

if [ -f "Makefile" ]; then
    echo "检测到 Makefile 文件"

    if command -v make &> /dev/null; then
        echo "正在检查 Makefile 目标..."

        if grep -q "^fmt:" Makefile; then
            echo "正在运行 make fmt..."
            make fmt 2>/dev/null && echo "✓ make fmt 完成" || echo "⚠ make fmt 失败或不存在"
        fi

    else
        echo "⚠ 未找到 make 命令，跳过 Makefile 操作"
    fi
else
    echo "提示: 未检测到 Makefile 文件"
fi

echo "=========================================="
echo "项目初始化完成！"
echo "=========================================="

