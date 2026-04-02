package apps

import (
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
	"git.bestfulfill.tech/devops/go-core/kits/kserver/gin_middles"
	"github.com/gin-gonic/gin"
)

// InitGinServer function 初始化 Gin 引擎.
// @autowire(set=init)
func InitGinServer() (eg *gin.Engine) {
	// 可按需要自行调整全局中间件
	eg = kserver.NewGinEngine()
	eg.Use(gin_middles.DefaultChain()...)
	return
}

// BaseServer struct 基础服务.
// @autowire(set=init)
type BaseServer struct {
	Runtime Runtime
	Engine  *gin.Engine
}

// Start method 启动服务.
// @config(cfg=ServerConfig)
func (server *BaseServer) Start(register func(router gin.IRouter), cfg *kserver.Config) {
	server.Runtime.Init()
	register(server.Engine)
	if cfg == nil {
		cfg = &kserver.Config{}
	}
	kserver.Run(server.Engine, *cfg)
}
