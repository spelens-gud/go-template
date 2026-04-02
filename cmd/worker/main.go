package main

import (
	"flag"
	"fmt"
	"os"

	"{{.ProjectName}}/cmd/internal"
	"{{.ProjectName}}/config"
	"{{.ProjectName}}/internal/apps"
	"{{.ProjectName}}/internal/version"

	"git.bestfulfill.tech/devops/go-core/kits/kenv/dotenv"
)

var _ = dotenv.ImportMe

func main() {
	showVersion := flag.Bool("v", false, "print version information")
	showVersionLong := flag.Bool("version", false, "print version information")

	flag.Parse()
	if *showVersion || *showVersionLong {
		fmt.Fprintln(os.Stdout, version.Get().String())
		return
	}

	apps.LoadConfigAndRun(func(c *config.Config) (app apps.App, cf func(), err error) {
		return internal.InitializeWorker(c)
	})
}
