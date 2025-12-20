package apps

import (
	"github.com/spelens-gud/Verktyg/kits/kgrpc"
	"github.com/spelens-gud/Verktyg/kits/kgrpc/interceptors"
	"google.golang.org/grpc"
)

// @autowire(set=init)
func InitGrpcServer() (gs *grpc.Server) {
	// 可按需要自行调整拦截器
	gs = kgrpc.NewGrpcServerWithInterceptors(
		interceptors.DefaultUnaryChain(),
		interceptors.DefaultStreamChain(),
	)
	return
}

// @autowire(set=init)
type BaseGrpcServer struct {
	Runtime    Runtime
	GrpcServer *grpc.Server
}

// Start method  启动 gRPC 服务器.
// @config(cfg=GrpcServerConfig)
func (server *BaseGrpcServer) Start(register func(gs *grpc.Server), cfg kgrpc.Config) {
	server.Runtime.Init()
	register(server.GrpcServer)
	kgrpc.Run(server.GrpcServer, cfg)
}
