package dao

import (
	"context"
	"go-template/internal/table"

	"gorm.io/gorm"
)

// @dao(OperationLogsDao)
// OperationLogsDao interface 操作日志接口.
type OperationLogsDao interface {
	// Create 创建操作日志.
	Create(ctx context.Context, db *gorm.DB, logisticInfo *table.OperationLogs) error
	// Update 更新操作日志.
	Update(ctx context.Context, db *gorm.DB, logisticInfo *table.OperationLogs) error
	// Delete 删除操作日志.
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	// GetByID 根据ID获取操作日志.
	GetByID(ctx context.Context, db *gorm.DB, id int64) (*table.OperationLogs, error)
	// List 获取操作日志列表.
	List(ctx context.Context, db *gorm.DB, query OperationLogsDaoQuery) ([]*table.OperationLogs, int64, error)
}

type OperationLogsDaoQuery struct {
	Name string `form:"name" json:"name"`
	PageList
}

type PageList struct {
	// 当前页码
	Page int `form:"page" json:"page"`
	// 每页数量
	PageSize int `form:"page_size" json:"page_size"`
	// 总数量
	Total int64 `form:"total" json:"total"`
}
