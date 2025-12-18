package dao_operation_logs_dao

import (
	"context"

	"gorm.io/gorm"

	"go-template/internal/table"
)

func (d *DaoImpl) Update(ctx context.Context, db *gorm.DB, logisticInfo *table.OperationLogs) error {
	panic("implement me")
}
