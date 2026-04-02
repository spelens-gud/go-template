package database

import (
	gormctx2 "git.bestfulfill.tech/devops/go-core/implements/gormctx/v2"
	"git.bestfulfill.tech/devops/go-core/interfaces/isql"
)

type DB isql.Gorm2SQL // gorm的db句柄

// InitSql method 初始化数据库连接.
// @autowire(set=db)
// @config(cfg=DBConfig)
func InitSql(cfg *isql.SQLConfig) (sql DB, cf func(), err error) {
	if sql, err = gormctx2.InitGormV2Sql(*cfg); err == nil {
		cf = func() {
			_ = sql.Close()
		}
	}
	return
}
