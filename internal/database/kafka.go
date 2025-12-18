package database

import (
	"github.com/Shopify/sarama"

	"git.bestfulfill.tech/devops/go-core/implements/skafka"
	"git.bestfulfill.tech/devops/go-core/interfaces/ikafka"
)

type (
	// 配置
	KafkaConfig ikafka.ClientConfig
	// 客户端
	KafkaClient ikafka.Client
	// 生产者
	KafkaProducer ikafka.Producer
)

// @autowire(set=db)
// @config(config)
// 初始化单点节点客户端
func InitKafkaClient(config KafkaConfig) (kafka KafkaClient, cf func(), err error) {
	if kafka, err = skafka.NewClient(config.Address, func(sc *sarama.Config) {
		sc.Version = sarama.V2_2_0_0
	}); err == nil {
		cf = func() { _ = kafka.Close() }
	}
	return
}

// @autowire(set=db)
// 初始化生产者
func InitKafkaProducer(kafka KafkaClient) (p KafkaProducer, cf func(), err error) {
	if p, err = kafka.NewAsyncProducer(); err == nil {
		cf = func() { _ = p.Close() }
	}
	return
}
