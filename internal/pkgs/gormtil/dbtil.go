package gormtil

import (
	"context"

	"gorm.io/gorm"
)

// Exist 判断表 table 中是否存在 condition 条件的记录
func Exist(db *gorm.DB, table, condition string) (exist bool, err error) {
	var count int64
	err = db.Table(table).Where(condition).Limit(1).Count(&count).Error
	if err != nil {
		return false, err
	}
	exist = count > 0
	return exist, nil
}

// DeleteByID 删除表 table 中 id 为 id 的记录
func DeleteByID(ctx context.Context, db *gorm.DB, table string, id int) (err error) {
	if table == "" {
		return
	}

	err = db.Table(table).Where("id = ?", id).Delete(nil).Error
	return
}
