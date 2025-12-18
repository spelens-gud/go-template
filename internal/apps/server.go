package apps

import (
	"github.com/gin-gonic/gin"

	"git.bestfulfill.tech/devops/go-core/kits/kserver"
	"git.bestfulfill.tech/devops/go-core/kits/kserver/gin_middles"
)

// @autowire(set=init)
func InitGinServer() (eg *gin.Engine) {
	// 可按需要自行调整全局中间件
	eg = kserver.NewGinEngine()
	eg.Use(gin_middles.DefaultChain()...)
	return
}

// @autowire(set=init)
type BaseServer struct {
	Runtime Runtime
	Engine  *gin.Engine
}

func (server *BaseServer) Start(register func(router gin.IRouter), cfg kserver.Config) {
	server.Runtime.Init()
	register(server.Engine)
	kserver.Run(server.Engine, cfg)
}
