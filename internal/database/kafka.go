package database

import (
	"git.bestfulfill.tech/devops/go-core/implements/skafka"
	"git.bestfulfill.tech/devops/go-core/interfaces/ikafka"
	"github.com/Shopify/sarama"
)

type (
	KafkaConfig   ikafka.ClientConfig // 配置
	KafkaClient   ikafka.Client       // 客户端
	KafkaProducer ikafka.Producer     // 生产者
)

// InitKafkaClient 初始化Kafka客户端.
// @autowire(set=db)
// InitKafkaClient 初始化单点节点客户端.
func InitKafkaClient(config KafkaConfig) (kafka KafkaClient, cf func(), err error) {
	if kafka, err = skafka.NewClient(config.Address, func(sc *sarama.Config) {
		sc.Version = sarama.V2_2_0_0
	}); err == nil {
		cf = func() { _ = kafka.Close() }
	}
	return
}

// InitKafkaProducer 创建Kafka生产者.
// @autowire(set=db)
// InitKafkaProducer 初始化生产者.
func InitKafkaProducer(kafka KafkaClient) (p KafkaProducer, cf func(), err error) {
	if p, err = kafka.NewAsyncProducer(); err == nil {
		cf = func() { _ = p.Close() }
	}
	return
}
