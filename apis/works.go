package apis

import (
	"github.com/spelens-gud/Verktyg/interfaces/iworker"
)

func (svc *Services) RegisterWorks(worker iworker.WorkerManager) {
	// 添加由context控制的阻塞任务 如队列消费等
	//worker.MustRegisterTask(svc.KafkaClient.Consume("test-group", []string{"test-topic"},
	//	func(ctx context.Context, message []byte) error {
	//		logger.FromContext(ctx).Infof("consume message: %s", message)
	//		return nil
	//	}))
	//
	//// 添加由crontab表达式控制的定时任务
	//worker.RegisterCron("0 * * * * *", func(ctx context.Context) (err error) {
	//	return svc.KafkaProducer.SendMessageX(ctx, "test-topic", "key", struct {
	//		Time time.Time
	//	}{time.Now()})
	//})
}
