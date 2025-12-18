package config

import (
	"go-template/internal/database"

	"git.bestfulfill.tech/devops/go-core/implements/otrace"
	"git.bestfulfill.tech/devops/go-core/implements/promgateway"
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
)

// @autowire.config()
// @mount(sql_config)
type Config struct {
	// 服务器配置
	ServerConfig kserver.Config `json:"server_config"`
	// 链路追踪配置
	TracerConfig otrace.JaegerConfig `json:"tracer_config"`
	// PrometheusGateway配置
	MetricsGatewayConfig promgateway.GatewayConfig `json:"metrics_gateway_config"`
	// 数据库配置
	MysqlConfig database.MysqlConfig `json:"db_config"`
}
