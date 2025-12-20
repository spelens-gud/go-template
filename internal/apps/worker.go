package apps

import (
	"github.com/robfig/cron/v3"
	"github.com/spelens-gud/Verktyg/implements/worker"
	"github.com/spelens-gud/Verktyg/interfaces/iworker"
)

// @autowire(set=init)
func InitWorkerManager() iworker.WorkerManager {
	return worker.NewWorkerManager(worker.WithCronOption(cron.WithSeconds()))
}

// @autowire(set=init)
type BaseWorker struct {
	Runtime       Runtime
	WorkerManager iworker.WorkerManager
}

func (worker *BaseWorker) Start(register func(manager iworker.WorkerManager)) {
	worker.Runtime.Init()
	register(worker.WorkerManager)
	worker.WorkerManager.Start()
}
