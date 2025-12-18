package dao_operation_logs_dao

import (
	"context"

	"gorm.io/gorm"

	dao "go-template/dao"
	"go-template/internal/table"
)

func (d *DaoImpl) List(ctx context.Context, db *gorm.DB, query dao.OperationLogsDaoQuery) ([]*table.OperationLogs, int64, error) {
	panic("implement me")
}
