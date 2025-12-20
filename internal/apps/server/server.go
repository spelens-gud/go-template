package server

import (
	"github.com/spelens-gud/Verktyg/kits/kgrpc"
	"google.golang.org/grpc"
	"{{.ProjectName}}/apis"

	"github.com/gin-gonic/gin"
	"github.com/spelens-gud/Verktyg/kits/kserver"

	"github.com/spelens-gud/Verktyg/interfaces/iconfig"
	"github.com/spelens-gud/Verktyg/kits/kstruct/structgraphx"

	"{{.ProjectName}}/internal/apps"
)

// @autowire.init()
type Server struct {
	BaseServer       apps.BaseServer     `json:"base_server"`
	BaseGrpcServer   apps.BaseGrpcServer `json:"base_grpc_server"`
	ServerConfig     kserver.Config      `json:"server_config"`
	GrpcServerConfig kgrpc.Config        `json:"grpc_server_config"`
	Services         apis.Services       `json:"services"`
	GrpcServices     apis.GrpcServices   `json:"grpc_services"`
}

// Run method  同时启动 Gin 和 gRPC 服务器.
func (app *Server) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_server.png")
	}

	// 启动 gRPC 服务器
	go app.BaseGrpcServer.Start(func(gs *grpc.Server) {
		app.GrpcServices.RegisterRouter(gs)
	}, app.GrpcServerConfig)

	// 启动 Gin 服务器
	app.BaseServer.Start(func(router gin.IRouter) {
		app.Services.RegisterRouter(router)
	}, app.ServerConfig)
}
