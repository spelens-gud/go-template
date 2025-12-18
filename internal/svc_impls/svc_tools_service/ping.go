package svc_tools_service

import (
	"context"
	"go-template/internal/table"
	service "go-template/service"
)

func (s *Service) Ping(ctx context.Context, params service.PingParam) (res service.PingRes, err error) {
	_ = s.OperationLogsDao.Create(ctx, s.SQL.GetDB(ctx), &table.OperationLogs{
		OperationType: "",
		Name:          "test",
		RecordId:      1,
		IpAddress:     "测试地址",
		UserAgent:     "admin",
	})

	if params.Data == "ping" {
		return service.PingRes{Data: "pong"}, nil
	}

	return
}
