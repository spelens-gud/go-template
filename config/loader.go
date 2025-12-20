package config

import (
	"github.com/spelens-gud/Verktyg/implements/cfgloader"
)

//var loader = cfgloader.NewConfigLoaderFromEnv()

var loader = cfgloader.NewFileLoader("./config/config.local.json")

func LoadConfig() (c Config) { loader.MustLoad(&c); return }
