package service

import (
	"context"
)

// ToolsService 日志服务.
// @service(tools,route="/tools")
type ToolsService interface {
	// @http(method=get,route="ping")
	Ping(ctx context.Context, params PingParam) (res PingRes, err error)
}

type PingParam struct {
	Data string `json:"data"`
}

type PingRes struct {
	Data string `json:"data"`
}
