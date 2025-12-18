package svc_tools_service

import (
	"context"

	"{{.ProjectName}}/internal/table"
	service "{{.ProjectName}}/service"
)

func (s *Service) Ping(ctx context.Context, params service.PingParam) (res string, err error) {
	_ = s.OperationLogsDao.Create(ctx, s.Db.GetDB(ctx), &table.OperationLogs{
		OperationType: "",
		Name:          "test",
		RecordId:      1,
		IpAddress:     "测试地址",
		UserAgent:     "admin",
	})

	if params.Data == "ping" {
		return "pong", nil
	}

	return
}
