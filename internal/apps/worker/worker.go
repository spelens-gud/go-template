package worker

import (
	"github.com/spelens-gud/Verktyg/interfaces/iconfig"
	"github.com/spelens-gud/Verktyg/interfaces/iworker"
	"github.com/spelens-gud/Verktyg/kits/kstruct/structgraphx"

	"{{.ProjectName}}/apis"
	"{{.ProjectName}}/internal/apps"
)

// @autowire.init()
type Worker struct {
	Services   apis.Services   `json:"services"`
	BaseWorker apps.BaseWorker `json:"base_worker"`
}

func (app *Worker) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_worker.png")
	}
	app.BaseWorker.Start(func(manager iworker.WorkerManager) {
		app.Services.RegisterWorks(manager)
	})
}
