package config

import (
	"{{.ProjectName}}/internal/database"

	"github.com/spelens-gud/Verktyg/implements/otrace"
	"github.com/spelens-gud/Verktyg/implements/promgateway"
)

// @autowire.config()
// @mount(db_config,server_config,config)
type Config struct {
	DbConfig             database.DBConfig         `json:"db_config"`
	MetricsGatewayConfig promgateway.GatewayConfig `json:"metrics_gateway_config"`
	TracerConfig         otrace.JaegerConfig       `json:"tracer_config"`
	ServerConfig         database.ServerConfig     `json:"server_config"`
}
