package database

import (
	gormctx2 "github.com/spelens-gud/Verktyg/implements/gormctx/v2"
	"github.com/spelens-gud/Verktyg/interfaces/isql"
)

// @db()
type DB isql.Gorm2SQL // gorm的db句柄

// @config()
type DBConfig isql.SQLConfig // 数据库配置

// @autowire(set=db)
func InitSql(config DBConfig) (sql DB, cf func(), err error) {
	if sql, err = gormctx2.InitGormV2Sql(isql.SQLConfig(config)); err == nil {
		cf = func() {
			_ = sql.Close()
		}
	}
	return
}
