package database

import (
	"git.bestfulfill.tech/devops/go-core/implements/redisx"
	"git.bestfulfill.tech/devops/go-core/interfaces/iredis"
)

type (
	Redis        iredis.Redis
	RedisCluster iredis.Redis
)

// InitRedis method 初始化主从 Redis
// @config(config=RedisConfig)
// @autowire(set=db)
func InitRedis(config *iredis.RedisConfig) (redis Redis, cf func(), err error) {
	if redis, err = redisx.NewRedis(config); err == nil {
		cf = func() { _ = redis.Close() }
	}
	return
}

// InitRedisCluster method 初始化集群 Redis
// 去掉x后就可以注入config配置依赖
// @autowire(set=db)
func InitRedisCluster(config *iredis.RedisClusterConfig) (redis RedisCluster, cf func(), err error) {
	if redis, err = redisx.NewRedisCluster(config); err == nil {
		cf = func() { _ = redis.Close() }
	}
	return
}
