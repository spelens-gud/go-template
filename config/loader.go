package config

import (
	"github.com/spelens-gud/Verktyg/implements/cfgloader"
)

var loader = cfgloader.NewConfigLoaderFromEnv()

func LoadConfig() (c Config) { loader.MustLoad(&c); return }
