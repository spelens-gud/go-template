package apis

import "{{.ProjectName}}/service"

// @autowire(set=service)
// @mount(service)
// @api_server()
type Services struct {
	ToolsService service.ToolsService `json:"tools_service"`
}
