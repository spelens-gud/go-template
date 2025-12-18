package apps

import (
	_ "go.uber.org/automaxprocs"

	"github.com/spelens-gud/Verktyg/implements/otrace"
	"github.com/spelens-gud/Verktyg/implements/promgateway"
	"github.com/spelens-gud/Verktyg/interfaces/imetrics"
	"github.com/spelens-gud/Verktyg/interfaces/itrace"
	"github.com/spelens-gud/Verktyg/kits/ktrace/tracerinit"
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
	return tracerinit.InitTracers(config)
}
