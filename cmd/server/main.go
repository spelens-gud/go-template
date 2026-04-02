package main

import (
	"flag"
	"fmt"
	"os"

	"{{.ProjectName}}/cmd/internal"
	"{{.ProjectName}}/config"
	_ "{{.ProjectName}}/docs"
	"{{.ProjectName}}/internal/apps"
	"{{.ProjectName}}/internal/version"

	"git.bestfulfill.tech/devops/go-core/kits/kenv/dotenv"
	"git.bestfulfill.tech/devops/go-core/kits/kserver/govern_server"
)

var _ = dotenv.ImportMe

// @title           系统api
// @version         1.0
// @description     系统服务
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:6066
// @BasePath  /api/v1
// @securityDefinitions.basic  BasicAuth
func main() {
	showVersion := flag.Bool("v", false, "print version information")
	showVersionLong := flag.Bool("version", false, "print version information")

	flag.Parse()
	if *showVersion || *showVersionLong {
		fmt.Fprintln(os.Stdout, version.Get().String())
		return
	}

	// 启动监控采集 server
	go govern_server.RunGovernServer(":9099")

	apps.LoadConfigAndRun(func(c *config.Config) (app apps.App, cf func(), err error) {
		return internal.InitializeServer(c)
	})
}
