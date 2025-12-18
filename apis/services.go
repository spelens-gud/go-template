package apis

import (
	service "go-template/service"
)

// @autowire(set=service)
// @mount(service)
type Services struct {
	ToolsService service.ToolsService
}
