package apis

import (
	"github.com/gin-gonic/gin"
)

const groupPrefix = "/api/opsclaw/v1"

func (svc *Services) RegisterRouter(group gin.IRouter) {
	// 路由前缀
	group = group.Group(groupPrefix)

	// TODO 统一控制器

	// TODO 中间价

	// TODO 路由注册
}
