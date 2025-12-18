#!/bin/bash

TOOLS=(
    "github.com/google/wire/cmd/wire@latest"
    "github.com/spf13/cobra-cli@latest"
    "golang.org/x/tools/cmd/goimports@latest"
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    "golang.org/x/vuln/cmd/govulncheck@latest"
)

if ! command -v go &> /dev/null; then
    echo "错误: 未找到 go 命令" >&2
    exit 1
fi

# 遍历工具列表并安装
for tool in "${TOOLS[@]}"; do
    tool_name=$(echo "$tool" | cut -d '@' -f1)
    if [[ "$tool_name" == *"/"* ]]; then
        binary_name=$(basename "$tool_name")
    else
        binary_name="$tool_name"
    fi

    if ! command -v "$binary_name" &> /dev/null; then
        if go install "$tool" &> /dev/null; then
            echo "✓ 安装 $binary_name"
        fi
    fi
done

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
            return 0
        fi
        check_name="typos"
    fi

    if command -v "$check_name" &> /dev/null || command -v "$tool_name" &> /dev/null; then
        return 0
    fi

    case "$os_type" in
        macos)
            if ! command -v brew &> /dev/null; then
                return 1
            fi
            case "$tool_name" in
                pre-commit)
                    brew install pre-commit &> /dev/null || pip3 install pre-commit &> /dev/null
                    ;;
                typos-cli)
                    brew install typos &> /dev/null || (command -v cargo &> /dev/null && cargo install typos-cli &> /dev/null)
                    ;;
                git-cliff)
                    brew install git-cliff &> /dev/null || (command -v cargo &> /dev/null && cargo install git-cliff &> /dev/null)
                    ;;
                make)
                    if ! command -v make &> /dev/null; then
                        xcode-select --install 2>/dev/null || true
                    fi
                    ;;
            esac
            ;;
        debian)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit &> /dev/null || ($sudo_cmd apt-get update &> /dev/null && $sudo_cmd apt-get install -y pre-commit &> /dev/null) || pip3 install pre-commit &> /dev/null
                    ;;
                typos-cli)
                    command -v cargo &> /dev/null && cargo install typos-cli &> /dev/null || return 1
                    ;;
                git-cliff)
                    command -v cargo &> /dev/null && cargo install git-cliff &> /dev/null || return 1
                    ;;
                make)
                    $sudo_cmd apt-get update &> /dev/null && $sudo_cmd apt-get install -y build-essential &> /dev/null
                    ;;
            esac
            ;;
        rhel|fedora)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit &> /dev/null || ([[ "$os_type" == "rhel" ]] && $sudo_cmd yum install -y pre-commit &> /dev/null || $sudo_cmd dnf install -y pre-commit &> /dev/null) || pip3 install pre-commit &> /dev/null
                    ;;
                typos-cli)
                    command -v cargo &> /dev/null && cargo install typos-cli &> /dev/null || return 1
                    ;;
                git-cliff)
                    command -v cargo &> /dev/null && cargo install git-cliff &> /dev/null || return 1
                    ;;
                make)
                    if [[ "$os_type" == "rhel" ]]; then
                        $sudo_cmd yum install -y make gcc &> /dev/null
                    else
                        $sudo_cmd dnf install -y make gcc &> /dev/null
                    fi
                    ;;
            esac
            ;;
        arch)
            case "$tool_name" in
                pre-commit)
                    pip3 install --user pre-commit &> /dev/null || $sudo_cmd pacman -S --noconfirm pre-commit &> /dev/null || pip3 install pre-commit &> /dev/null
                    ;;
                typos-cli)
                    $sudo_cmd pacman -S --noconfirm typos &> /dev/null || (command -v cargo &> /dev/null && cargo install typos-cli &> /dev/null) || true
                    ;;
                git-cliff)
                    $sudo_cmd pacman -S --noconfirm git-cliff &> /dev/null || (command -v cargo &> /dev/null && cargo install git-cliff &> /dev/null) || true
                    ;;
                make)
                    $sudo_cmd pacman -S --noconfirm make &> /dev/null
                    ;;
            esac
            ;;
        *)
            return 1
            ;;
    esac

    if [[ "$tool_name" == "typos-cli" ]]; then
        command -v typos &> /dev/null && echo "✓ 安装 $tool_name" && return 0
    fi

    if command -v "$tool_name" &> /dev/null; then
        echo "✓ 安装 $tool_name"
        return 0
    fi
    return 1
}

OS_TYPE=$(detect_os)
SYSTEM_TOOLS=("pre-commit" "typos-cli" "git-cliff" "make")

for tool in "${SYSTEM_TOOLS[@]}"; do
    install_tool "$tool" "$OS_TYPE" &> /dev/null
done
