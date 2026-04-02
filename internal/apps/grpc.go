package apps

import (
	"context"

	"git.bestfulfill.tech/devops/go-core/kits/kgrpc"
	"git.bestfulfill.tech/devops/go-core/kits/kgrpc/interceptors"
	"google.golang.org/grpc"
)

// InitGrpcServer function 初始化 gRPC 服务器.
// @autowire(set=init)
func InitGrpcServer() (gs *grpc.Server) {
	// 可按需要自行调整拦截器
	gs = kgrpc.NewGrpcServerWithInterceptors(
		interceptors.DefaultUnaryChain(),
		interceptors.DefaultStreamChain(),
	)
	return
}

// BaseGrpcServer struct gRPC 服务器基类.
// @autowire(set=init)
type BaseGrpcServer struct {
	Runtime    Runtime
	GrpcServer *grpc.Server
}

// Start method 启动 gRPC 服务器.
// @config(cfg=GrpcServerConfig)
func (server *BaseGrpcServer) Start(register func(gs *grpc.Server), cfg *kgrpc.ServerConfig) {
	if cfg == nil {
		cfg = &kgrpc.ServerConfig{}
	}
	server.Runtime.Init()
	register(server.GrpcServer)
	kgrpc.RunWithNacosServerContext(context.Background(), server.GrpcServer, *cfg)
}
