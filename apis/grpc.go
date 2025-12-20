package apis

import (
	"google.golang.org/grpc"
	"{{.ProjectName}}/internal/grpc_impls/grpc_user_server"
	"{{.ProjectName}}/proto"
)

// @autowire(set=grpc)
type GrpcServices struct {
	UserServer *grpc_user_server.GrpcImpl
}

func (svc *GrpcServices) RegisterRouter(gs *grpc.Server) {
	// 注册grpc服务
	proto.RegisterUserServer(gs, svc.UserServer)
}
