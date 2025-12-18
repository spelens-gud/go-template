package worker

import (
	"git.bestfulfill.tech/devops/go-core/interfaces/iconfig"
	"git.bestfulfill.tech/devops/go-core/interfaces/iworker"
	"git.bestfulfill.tech/devops/go-core/kits/kstruct/structgraphx"

	"go-template/apis"
	"go-template/internal/apps"
)

// @autowire.init()
type Worker struct {
	Services   apis.Services
	BaseWorker apps.BaseWorker
}

func (app *Worker) Run() {
	if iconfig.GetEnv().IsDevelopment() {
		go structgraphx.GenStructGraph(app, "design/structure_worker.png")
	}
	app.BaseWorker.Start(func(manager iworker.WorkerManager) {
		app.Services.RegisterWorks(manager)
	})
}
