package config

import (
	"git.bestfulfill.tech/devops/go-core/implements/cfgloader"
)

var loader = cfgloader.NewFileLoader("./config/config.json")

//var loader = cfgloader.NewConfigLoaderFromEnv()

func LoadConfig() (c Config) { loader.MustLoad(&c); return }
