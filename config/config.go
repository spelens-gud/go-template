package config

import (
	"github.com/spelens-gud/Verktyg/implements/otrace"
	"github.com/spelens-gud/Verktyg/implements/promgateway"
	"github.com/spelens-gud/Verktyg/kits/kgrpc"
	"github.com/spelens-gud/Verktyg/kits/kserver"

	"{{.ProjectName}}/internal/database"
)

// @autowire.config()
// @mount(config)
type Config struct {
	GrpcServerConfig     kgrpc.Config              `json:"grpc_server_config"`
	ServerConfig         kserver.Config            `json:"server_config"`
	MetricsGatewayConfig promgateway.GatewayConfig `json:"metrics_gateway_config"`
	TracerConfig         otrace.JaegerConfig       `json:"tracer_config"`
	DbConfig             database.DBConfig         `json:"db_config"`
}
