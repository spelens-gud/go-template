package grpc_user_server

import (
	"context"

	"{{.ProjectName}}/proto"
)

func (g *GrpcImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	panic("implement me")
}
