package apps

import (
	_ "go.uber.org/automaxprocs"

	"git.bestfulfill.tech/devops/go-core/implements/otrace"
	"git.bestfulfill.tech/devops/go-core/implements/promgateway"
	"git.bestfulfill.tech/devops/go-core/interfaces/imetrics"
	"git.bestfulfill.tech/devops/go-core/interfaces/itrace"
	"git.bestfulfill.tech/devops/go-core/kits/ktrace/tracerinit"
)

// @autowire(set=init)
type Runtime struct {
	MetricsPushDaemon imetrics.GatewayDaemon `structgraph:"-"`
	Tracer            itrace.Tracer          `structgraph:"-"`
}

func (t *Runtime) Init() {}

// @autowire(set=init)
// @config(gatewayCfg=MetricsGatewayConfig)
func InitMetricsPush(gatewayCfg promgateway.GatewayConfig) (daemon imetrics.GatewayDaemon, cf func()) {
	daemon = gatewayCfg.NewTransport()
	go daemon.StartDaemon()
	return daemon, daemon.Stop
}

// @autowire(set=init)
// @config(config=TracerConfig)
func InitTracer(config otrace.JaegerConfig) (tc itrace.Tracer, cleanup func(), err error) {
	return tracerinit.InitTracer(config)
}
