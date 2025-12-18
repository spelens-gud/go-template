package dao_operation_logs_dao

import (
	"context"

	"gorm.io/gorm"

	"{{.ProjectName}}/internal/table"
)

func (d *DaoImpl) GetByID(ctx context.Context, db *gorm.DB, id int64) (*table.OperationLogs, error) {
	panic("implement me")
}
