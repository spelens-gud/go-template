package server

import (
	"{{.ProjectName}}/apis"
	"{{.ProjectName}}/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/spelens-gud/Verktyg/kits/kserver"

	"github.com/spelens-gud/Verktyg/interfaces/iconfig"
	"github.com/spelens-gud/Verktyg/kits/kstruct/structgraphx"

	"{{.ProjectName}}/internal/apps"
)

// @autowire.init()
// @mount(api_server,base_server,server_config)
type Server struct {
	BaseServer   apps.BaseServer       `json:"base_server"`
	ServerConfig database.ServerConfig `json:"server_config"`
	Services     apis.Services         `json:"services"`
}

func (app *Server) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_server.png")
	}
	app.BaseServer.Start(func(router gin.IRouter) {
		app.Services.RegisterRouter(router)
	}, kserver.Config(app.ServerConfig))
}
