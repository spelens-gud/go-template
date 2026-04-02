package apis

import (
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
	"github.com/gin-gonic/gin"
)

// 默认使用标准控制器 可根据业务需求定制
var defaultController = &kserver.DefaultServiceHandler{}

func (svc *Services) Response(c *gin.Context, ret interface{}, err error) {
	defaultController.Response(c, ret, err)
}

func (svc *Services) ParamHandler(c *gin.Context, params []interface{}) bool {
	return defaultController.ParamHandler(c, params)
}
