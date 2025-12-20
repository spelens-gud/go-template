package grpc_user_server

import (
	"context"

	"{{.ProjectName}}/proto"
)

func (g *GrpcImpl) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	panic("implement me")
}
