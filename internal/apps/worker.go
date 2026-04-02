package apps

import (
	"git.bestfulfill.tech/devops/go-core/implements/worker"
	"git.bestfulfill.tech/devops/go-core/interfaces/iworker"
)

// InitWorkerManager function 初始化 WorkerManager.
// @autowire(set=init)
// @config(cfg=WorkerConfig)
func InitWorkerManager(cfg *worker.Config) iworker.WorkerManager {
	if cfg == nil {
		cfg = &worker.Config{}
	}
	return worker.NewWorkerManagerFromConfig(*cfg)
}

// BaseWorker struct 基础 Worker.
// @autowire(set=init)
type BaseWorker struct {
	Runtime       Runtime
	WorkerManager iworker.WorkerManager
}

func (w *BaseWorker) Start(register func(manager iworker.WorkerManager)) {
	w.Runtime.Init()
	register(w.WorkerManager)
	w.WorkerManager.Start()
}
