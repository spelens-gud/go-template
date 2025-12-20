package grpc_user_server

import (
	"context"

	"{{.ProjectName}}/proto"
)

func (g *GrpcImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	panic("implement me")
}
