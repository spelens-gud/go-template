package dao_operation_logs_dao

import (
	"context"

	"gorm.io/gorm"
)

func (d *DaoImpl) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	panic("implement me")
}
