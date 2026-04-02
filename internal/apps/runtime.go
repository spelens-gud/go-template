package apps

import (
	"git.bestfulfill.tech/devops/go-core/implements/otrace"
	"git.bestfulfill.tech/devops/go-core/implements/promgateway"
	"git.bestfulfill.tech/devops/go-core/interfaces/imetrics"
	"git.bestfulfill.tech/devops/go-core/interfaces/itrace"
	"git.bestfulfill.tech/devops/go-core/kits/ktrace/tracerinit"
	_ "go.uber.org/automaxprocs"
)

// Runtime struct 运行时环境基类.
// @autowire(set=init)
type Runtime struct {
	MetricsPushDaemon imetrics.GatewayDaemon `structgraph:"-"`
	Tracer            itrace.Tracer          `structgraph:"-"`
}

func (t *Runtime) Init() {}

// InitMetricsPush function 初始化 Metrics 推送服务.
// @autowire(set=init)
// @config(gatewayCfg=MetricsGatewayConfig)
func InitMetricsPush(gatewayCfg *promgateway.GatewayConfig) (daemon imetrics.GatewayDaemon, cf func()) {
	if gatewayCfg == nil {
		gatewayCfg = &promgateway.GatewayConfig{}
	}
	daemon = gatewayCfg.NewTransport()
	go daemon.StartDaemon()
	return daemon, daemon.Stop
}

// InitTracer function 创建 Jaeger 追踪器.
// @autowire(set=init)
// @config(config=TracerConfig)
func InitTracer(config *otrace.JaegerConfig) (tc itrace.Tracer, cleanup func(), err error) {
	if config == nil {
		config = &otrace.JaegerConfig{}
	}
	return tracerinit.InitTracers(config)
}
