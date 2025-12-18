package database

import (
	"git.bestfulfill.tech/devops/go-core/implements/redisx"
	"git.bestfulfill.tech/devops/go-core/interfaces/iredis"
)

type (
	// 主从节点
	Redis       iredis.Redis
	RedisConfig *iredis.RedisConfig

	// 集群
	RedisCluster       iredis.Redis
	RedisClusterConfig *iredis.RedisClusterConfig
)

// @autowire(set=db)
// @config.x(config)
// 初始化主从
func InitRedis(config RedisConfig) (redis Redis, cf func(), err error) {
	if redis, err = redisx.NewRedis(config); err == nil {
		cf = func() { _ = redis.Close() }
	}
	return
}

// @autowire(set=db)
// @config.x(config)
// 初始化集群
func InitRedisCluster(config RedisClusterConfig) (redis RedisCluster, cf func(), err error) {
	if redis, err = redisx.NewRedisCluster(config); err == nil {
		cf = func() { _ = redis.Close() }
	}
	return
}
