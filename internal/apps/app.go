package apps

import (
	"fmt"

	"go-template/config"
)

type App interface {
	Run()
}

func LoadConfigAndRun(init func(*config.Config) (app App, cf func(), err error)) {
	cfg := config.LoadConfig()
	app, cleanup, err := init(&cfg)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	defer cleanup()
	app.Run()
}
