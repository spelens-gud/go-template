package apis

import (
	"google.golang.org/grpc"
)

// GrpcServices struct gRPC服务.
// @autowire(set=grpc)
type GrpcServices struct {
	//UserServer *grpc_user_server.GrpcImpl
}

func (svc *GrpcServices) RegisterRouter(gs *grpc.Server) {
	// 注册grpc服务
	//proto.RegisterUserServer(gs, svc.UserServer)
}
