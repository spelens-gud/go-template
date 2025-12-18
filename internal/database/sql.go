package database

import (
	gormctx2 "git.bestfulfill.tech/devops/go-core/implements/gormctx/v2"
	"git.bestfulfill.tech/devops/go-core/interfaces/isql"
)

type (
	Mysql       isql.Gorm2SQL
	MysqlConfig isql.SQLConfig
)

// 单个sql数据库
// @autowire(set=db)
// @sql_config(config=MysqlConfig)
func InitSql(config MysqlConfig) (sql Mysql, cf func(), err error) {
	if sql, err = gormctx2.InitGormV2Sql(isql.SQLConfig(config)); err == nil {
		cf = func() {
			_ = sql.Close()
		}
	}
	return
}
