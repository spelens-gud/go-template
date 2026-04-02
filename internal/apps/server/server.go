package server

import (
	"context"

	"{{.ProjectName}}/apis"
	"{{.ProjectName}}/internal/apps"
	"{{.ProjectName}}/internal/pkgs/gotil"

	"git.bestfulfill.tech/devops/go-core/interfaces/iconfig"
	"git.bestfulfill.tech/devops/go-core/interfaces/iworker"
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
	"git.bestfulfill.tech/devops/go-core/kits/kstruct/structgraphx"
	"github.com/gin-gonic/gin"
)

// Server struct http服务.
// @autowire.init()
type Server struct {
	BaseServer apps.BaseServer `json:"base_server"`
	BaseWorker apps.BaseWorker `json:"base_worker"`
	// BaseGrpcServer   apps.BaseGrpcServer `json:"base_grpc_server"`
	ServerConfig *kserver.Config `json:"server_config"`
	// GrpcServerConfig kgrpc.ServerConfig  `json:"grpc_server_config"`
	Services     apis.Services     `json:"services"`
	WorkServices apis.WorkServices `json:"work_services"`
	// GrpcServices     apis.GrpcServices   `json:"grpc_services"`
}

// Run method  同时启动 Gin 和 gRPC 服务器.
func (app *Server) Run() {
	if err := gotil.GoErrGroup(context.Background(),
		func(ctx context.Context) (err error) {
			if iconfig.GetEnv().IsDevelopment() {
				structgraphx.GenStructGraph(app, "design/structure_server.png")
			}
			return nil
		},
		// 启动定时服务
		func(ctx context.Context) (err error) {
			app.BaseWorker.Start(func(manager iworker.WorkerManager) {
				app.WorkServices.RegisterWorks(manager)
			})
			return nil
		},
		// 启动 gRPC 服务器
		func(ctx context.Context) (err error) {
			// TODO 启动grpc
			// app.BaseGrpcServer.Start(func(gs *grpc.Server) {
			//	app.GrpcServices.RegisterRouter(gs)
			// }, app.GrpcServerConfig)
			return nil
		},
		// 启动 Gin 服务器
		func(ctx context.Context) (err error) {
			if app.ServerConfig == nil {
				app.ServerConfig = &kserver.Config{}
			}
			app.BaseServer.Start(func(router gin.IRouter) {
				app.Services.RegisterRouter(router)
			}, app.ServerConfig)
			return nil
		},
	); err != nil {
		return
	}
}
