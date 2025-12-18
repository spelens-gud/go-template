package dao_operation_logs_dao

import (
	"context"

	"gorm.io/gorm"

	"{{.ProjectName}}/internal/table"
)

func (d *DaoImpl) Create(ctx context.Context, db *gorm.DB, logisticInfo *table.OperationLogs) error {
	return db.Create(logisticInfo).Error
}
