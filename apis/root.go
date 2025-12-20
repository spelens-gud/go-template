package apis

import (
	svrlessgin "github.com/Just-maple/serverless-gin"
	"github.com/gin-gonic/gin"
	"{{.ProjectName}}/apis/tools"
)

const groupPrefix = "/api/v1"

// @title API系统
// @version 1.0
// @description API系统接口
// @BasePath /api/v1
func (svc *Services) RegisterRouter(group gin.IRouter) {
	// 路由前缀
	group = group.Group(groupPrefix)

	// 统一控制器
	serverlessController := svrlessgin.NewWithController(svc)

	// 路由挂载...
	tools.RegisterToolsGroup(svc.ToolsService, group, serverlessController)
}
