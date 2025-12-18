package server

import (
	"go-template/apis"

	"github.com/gin-gonic/gin"

	"git.bestfulfill.tech/devops/go-core/interfaces/iconfig"
	"git.bestfulfill.tech/devops/go-core/kits/kserver"
	"git.bestfulfill.tech/devops/go-core/kits/kstruct/structgraphx"

	"go-template/internal/apps"
)

// @autowire.init()
type Server struct {
	Services   apis.Services
	BaseServer apps.BaseServer
	Config     kserver.Config `structgraph:"-"`
}

func (app *Server) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_server.png")
	}
	app.BaseServer.Start(func(router gin.IRouter) {
		app.Services.RegisterRouter(router)
	}, app.Config)
}
