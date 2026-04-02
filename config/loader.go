package config

import (
	"git.bestfulfill.tech/devops/go-core/implements/cfgloader"
)

var loader = cfgloader.NewMultiLoader(
	cfgloader.NewFileLoader("/acm-config-data/{{.FolderName}}"),
	// cfgloader.NewNacosLoader(cfgloader.WithOfficialSDK("./runtime")),
	// cfgloader.NewConfigLoaderFromEnv(),
	cfgloader.NewFileLoader("./config/config.local.json"),
)

func LoadConfig() (c Config) { loader.MustLoad(&c); return }
