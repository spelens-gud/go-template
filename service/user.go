package service

import (
	"context"

	"{{.ProjectName}}/proto"
)

// @grpc(user,proto="user")
type UserServer interface {
	// GetUser 获取用户信息
	GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error)
}
