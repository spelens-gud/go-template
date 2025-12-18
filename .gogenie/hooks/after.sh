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
/*
Copyright © 2025 xusihong

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"$PROJECT_NAME/internal/version"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "$CMD_NAME",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// defaultVersionTemplate 默认版本信息模板.
const defaultVersionTemplate = `{{with .DisplayName}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}

`

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// Execute function    执行根命令.
func Execute() {
	if term.IsTerminal(os.Stdout.Fd()) {
		var b bytes.Buffer
		w := colorprofile.NewWriter(os.Stdout, os.Environ())
		w.Forward = &b
		rootCmd.SetVersionTemplate(b.String() + "\n" + defaultVersionTemplate)
	}

	// 检测终端尺寸，用于优化输出格式
	var width, height int
	if term.IsTerminal(os.Stdout.Fd()) {
		var err error
		width, height, err = term.GetSize(os.Stdout.Fd())
		if err == nil && width > 0 {
			// 可以根据终端宽度调整输出格式
			_ = width
			_ = height
		}
	}

	// 执行命令，使用 fang 的增强功能
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version.Version),
		fang.WithNotifySignal(os.Interrupt),
		fang.WithErrorHandler(customErrorHandler),
	); err != nil {
		log.Error("command execution failed:", err.Error())
		os.Exit(1)
	}
}

// customErrorHandler 自定义错误处理器，使用颜色输出错误信息.
func customErrorHandler(w io.Writer, _ fang.Styles, err error) {
	writer := colorprofile.NewWriter(w, os.Environ())
	const colorRed = "\x1b[31m"
	const colorReset = "\x1b[0m"
	//nolint:errcheck
	fmt.Fprintf(writer, "%s错误: %v%s\n", colorRed, err, colorReset)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/$CMD_NAME.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
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
fi

if [ -f ".pre-commit-config.yaml" ] && command -v pre-commit &> /dev/null; then
    pre-commit install &> /dev/null && echo "✓ 安装 pre-commit hooks"
fi

if [ -f "Makefile" ] && command -v make &> /dev/null && grep -q "^fmt:" Makefile; then
    make fmt &> /dev/null && echo "✓ 格式化代码"
fi
