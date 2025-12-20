package grpc

import (
	"github.com/spelens-gud/Verktyg/interfaces/iconfig"
	"github.com/spelens-gud/Verktyg/kits/kgrpc"
	"github.com/spelens-gud/Verktyg/kits/kstruct/structgraphx"
	"google.golang.org/grpc"
	"{{.ProjectName}}/apis"
	"{{.ProjectName}}/internal/apps"
)

// @autowire.init()
type GrpcServer struct {
	BaseGrpcServer   apps.BaseGrpcServer `json:"base_grpc_server"`
	GrpcServerConfig kgrpc.Config        `json:"grpc_server_config"`
	Services         apis.GrpcServices   `json:"services"`
}

// Run method  启动 gRPC 服务器.
func (app *GrpcServer) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_grpc.png")
	}
	app.BaseGrpcServer.Start(func(gs *grpc.Server) {
		app.Services.RegisterRouter(gs)
	}, app.GrpcServerConfig)
}
