/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spelens-gud/Verktyg/interfaces/iconfig"
	"github.com/spelens-gud/Verktyg/kits/klog/logger"
	"github.com/spelens-gud/Verktyg/kits/kserver/govern_server"
	"github.com/spf13/cobra"
	"{{.ProjectName}}/cmd/internal"
	"{{.ProjectName}}/config"
	"{{.ProjectName}}/internal/apps"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		iconfig.SetEnv("Development")
		// 注册业务域/命令空间，后面上k8s 后可从环境变量中获取
		iconfig.SetApplicationNameSpace("oms")
		// 注册服务名，后面上k8s 后可从环境变量中获取
		iconfig.SetApplicationAppName("purchase")
		// 设置日志文件输出目录
		_ = logger.SetFileOutput(func(option *logger.FileLogOption) {
			option.Dir = "./logs"
		})

		// 启动监控采集server
		go govern_server.RunGovernServer(":9099")
		apps.LoadConfigAndRun(func(c *config.Config) (app apps.App, cf func(), err error) {
			return internal.InitializeGrpcServer(c)
		})
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
