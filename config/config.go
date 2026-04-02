package config

import (
	"git.bestfulfill.tech/devops/go-core/implements/otrace"
	"git.bestfulfill.tech/devops/go-core/implements/promgateway"
	"git.bestfulfill.tech/devops/go-core/implements/worker"
	"git.bestfulfill.tech/devops/go-core/interfaces/iredis"
	"git.bestfulfill.tech/devops/go-core/interfaces/isql"
	"git.bestfulfill.tech/devops/go-core/kits/kgrpc"
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
)

// Config struct 配置信息.
// @autowire.config()
// @mount(config)
type Config struct {
	GrpcServerConfig     *kgrpc.ServerConfig        `json:"grpc_server_config"`
	MetricsGatewayConfig *promgateway.GatewayConfig `json:"metrics_gateway_config"`
	TracerConfig         *otrace.JaegerConfig       `json:"tracer_config"`
	ServerConfig         *kserver.Config            `json:"server_config"`
	WorkerConfig         *worker.Config             `json:"worker_config"`
	RedisConfig          *iredis.RedisConfig        `json:"redis_config"`
	DbConfig             *isql.SQLConfig            `json:"db_config"`
}
