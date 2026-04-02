package grpc

import (
	"{{.ProjectName}}/apis"
	"{{.ProjectName}}/internal/apps"

	"git.bestfulfill.tech/devops/go-core/interfaces/iconfig"
	"git.bestfulfill.tech/devops/go-core/kits/kgrpc"
	"git.bestfulfill.tech/devops/go-core/kits/kstruct/structgraphx"
	"google.golang.org/grpc"
)

// GrpcServer struct gRPC 服务器.
// @autowire.init()
type GrpcServer struct {
	GrpcServerConfig *kgrpc.ServerConfig `json:"grpc_server_config"`
	BaseGrpcServer   apps.BaseGrpcServer `json:"base_grpc_server"`
	Services         apis.GrpcServices   `json:"services"`
}

// Run method 启动 gRPC 服务器.
func (app *GrpcServer) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_grpc.png")
	}
	if app.GrpcServerConfig == nil {
		app.GrpcServerConfig = &kgrpc.ServerConfig{}
	}
	app.BaseGrpcServer.Start(func(gs *grpc.Server) {
		app.Services.RegisterRouter(gs)
	}, app.GrpcServerConfig)
}
