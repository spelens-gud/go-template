#!/bin/bash

TOOLS=(
    "github.com/google/wire/cmd/wire@latest"
    "github.com/spf13/cobra-cli@latest"
    "golang.org/x/tools/cmd/goimports@latest"
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    "golang.org/x/vuln/cmd/govulncheck@latest"
)

if ! command -v go &> /dev/null; then
    echo "错误: 未找到 go 命令，请先安装 Go 环境"
    exit 1
fi

echo "当前 Go 版本: $(go version)"

# 遍历工具列表并安装
for tool in "${TOOLS[@]}"; do
    tool_name=$(echo "$tool" | cut -d '@' -f1)
    if [[ "$tool_name" == *"/"* ]]; then
        binary_name=$(basename "$tool_name")
    else
        binary_name="$tool_name"
    fi

    # 检查工具是否已经安装
    if command -v "$binary_name" &> /dev/null; then
        echo "✓ $binary_name 已经安装, 跳过"
    else
        echo "✗ $binary_name 未安装，正在安装..."
        if go install "$tool"; then
            echo "✓ 成功安装 $binary_name"
        else
            echo "✗ 安装 $binary_name 失败"
        fi
    fi
done

echo ""

detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # 检测 Linux 发行版
        if command -v apt-get &> /dev/null; then
            echo "debian"
        elif command -v yum &> /dev/null; then
            echo "rhel"
        elif command -v dnf &> /dev/null; then
            echo "fedora"
        elif command -v pacman &> /dev/null; then
            echo "arch"
        else
            echo "linux"
        fi
    else
        echo "unknown"
    fi
}

check_sudo() {
    if [[ $EUID -ne 0 ]]; then
        if command -v sudo &> /dev/null; then
            echo "sudo"
        else
            echo ""
        fi
    else
        echo ""
    fi
}

install_tool() {
    local tool_name=$1
    local os_type=$2
    local sudo_cmd=$(check_sudo)

    # 检查工具是否已安装（typos-cli 的命令名可能是 typos）
    local check_name="$tool_name"
    if [[ "$tool_name" == "typos-cli" ]]; then
        if command -v typos &> /dev/null; then
            echo "✓ typos (typos-cli) 已经安装, 跳过"
            return 0
        fi
        check_name="typos"
    fi

    if command -v "$check_name" &> /dev/null || command -v "$tool_name" &> /dev/null; then
        echo "✓ $tool_name 已经安装, 跳过"
        return 0
    fi

    echo "✗ $tool_name 未安装，正在安装..."

    case "$os_type" in
        macos)
            if ! command -v brew &> /dev/null; then
                echo "错误: 未找到 Homebrew，请先安装 Homebrew (https://brew.sh)"
                return 1
            fi
            case "$tool_name" in
                pre-commit)
                    brew install pre-commit || pip3 install pre-commit
                    ;;
                typos-cli)
                    brew install typos || (command -v cargo &> /dev/null && cargo install typos-cli)
                    ;;
                git-cliff)
                    brew install git-cliff || (command -v cargo &> /dev/null && cargo install git-cliff)
                    ;;
                make)
                    # macOS 通常自带 make
                    if ! command -v make &> /dev/null; then
                        echo "提示: 正在安装 Xcode Command Line Tools..."
                        xcode-select --install 2>/dev/null || echo "提示: 请手动安装 Xcode Command Line Tools"
                    fi
                    ;;
            esac
            ;;
        debian)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit || $sudo_cmd apt-get update && $sudo_cmd apt-get install -y pre-commit || pip3 install pre-commit
                    ;;
                typos-cli)
                    if command -v cargo &> /dev/null; then
                        cargo install typos-cli
                    else
                        echo "提示: typos-cli 需要 cargo，请先安装 Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
                        return 1
                    fi
                    ;;
                git-cliff)
                    if command -v cargo &> /dev/null; then
                        cargo install git-cliff
                    else
                        echo "提示: git-cliff 需要 cargo，请先安装 Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
                        return 1
                    fi
                    ;;
                make)
                    $sudo_cmd apt-get update && $sudo_cmd apt-get install -y build-essential
                    ;;
            esac
            ;;
        rhel|fedora)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit || ([[ "$os_type" == "rhel" ]] && $sudo_cmd yum install -y pre-commit || $sudo_cmd dnf install -y pre-commit) || pip3 install pre-commit
                    ;;
                typos-cli)
                    if command -v cargo &> /dev/null; then
                        cargo install typos-cli
                    else
                        echo "提示: typos-cli 需要 cargo，请先安装 Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
                        return 1
                    fi
                    ;;
                git-cliff)
                    if command -v cargo &> /dev/null; then
                        cargo install git-cliff
                    else
                        echo "提示: git-cliff 需要 cargo，请先安装 Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
                        return 1
                    fi
                    ;;
                make)
                    if [[ "$os_type" == "rhel" ]]; then
                        $sudo_cmd yum install -y make gcc
                    else
                        $sudo_cmd dnf install -y make gcc
                    fi
                    ;;
            esac
            ;;
        arch)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit || $sudo_cmd pacman -S --noconfirm pre-commit || pip3 install pre-commit
                    ;;
                typos-cli)
                    $sudo_cmd pacman -S --noconfirm typos 2>/dev/null || (command -v cargo &> /dev/null && cargo install typos-cli) || echo "提示: 请手动安装 typos-cli"
                    ;;
                git-cliff)
                    $sudo_cmd pacman -S --noconfirm git-cliff 2>/dev/null || (command -v cargo &> /dev/null && cargo install git-cliff) || echo "提示: 请手动安装 git-cliff"
                    ;;
                make)
                    $sudo_cmd pacman -S --noconfirm make
                    ;;
            esac
            ;;
        *)
            echo "警告: 未识别的操作系统类型，请手动安装 $tool_name"
            return 1
            ;;
    esac

    if [[ "$tool_name" == "typos-cli" ]]; then
        if command -v typos &> /dev/null; then
            echo "✓ 成功安装 typos (typos-cli)"
            return 0
        fi
    fi

    if command -v "$tool_name" &> /dev/null; then
        echo "✓ 成功安装 $tool_name"
        return 0
    else
        echo "✗ 安装 $tool_name 失败，请手动安装"
        return 1
    fi
}

OS_TYPE=$(detect_os)
echo "检测到操作系统类型: $OS_TYPE"

SYSTEM_TOOLS=("pre-commit" "typos-cli" "git-cliff" "make")

for tool in "${SYSTEM_TOOLS[@]}"; do
    install_tool "$tool" "$OS_TYPE"
done

echo ""
echo "所有依赖检查完成！"
