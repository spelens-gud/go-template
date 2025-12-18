package service

import (
	"context"
)

// ToolsService 日志服务.
// @service(tools,route="/tools")
type ToolsService interface {
	// @http(method=get,route="ping")
	Ping(ctx context.Context, params PingParam) (res string, err error)
}

type PingParam struct {
	Data string `json:"data"`
}
