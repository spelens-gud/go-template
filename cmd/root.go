// Package cmd provides CLI commands.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"go-template/internal/version"
	"io"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

const ProjectName = "go-template"

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   ProjectName,
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely
contains examples and usage of using your application. For example:

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
