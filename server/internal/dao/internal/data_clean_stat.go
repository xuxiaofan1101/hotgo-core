// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataCleanStatDao is the data access object for the table hg_data_clean_stat.
type DataCleanStatDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  DataCleanStatColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// DataCleanStatColumns defines and stores column names for the table hg_data_clean_stat.
type DataCleanStatColumns struct {
	Id               string // 清洗统计ID
	TaskId           string // 清洗任务ID
	ConnectorId      string // 输入连接ID
	EventType        string // 事件类型
	WindowStart      string // 统计窗口开始时间
	WindowEnd        string // 统计窗口结束时间
	TotalCount       string // 总处理数
	SuccessCount     string // 成功数
	FailedCount      string // 失败数
	CleanFailedCount string // 清洗失败数
	DroppedCount     string // 丢弃数
	LastError        string // 最近错误
	ErrorSamples     string // 错误样例
	CreatedAt        string // 创建时间
	UpdatedAt        string // 修改时间
}

// dataCleanStatColumns holds the columns for the table hg_data_clean_stat.
var dataCleanStatColumns = DataCleanStatColumns{
	Id:               "id",
	TaskId:           "task_id",
	ConnectorId:      "connector_id",
	EventType:        "event_type",
	WindowStart:      "window_start",
	WindowEnd:        "window_end",
	TotalCount:       "total_count",
	SuccessCount:     "success_count",
	FailedCount:      "failed_count",
	CleanFailedCount: "clean_failed_count",
	DroppedCount:     "dropped_count",
	LastError:        "last_error",
	ErrorSamples:     "error_samples",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewDataCleanStatDao creates and returns a new DAO object for table data access.
func NewDataCleanStatDao(handlers ...gdb.ModelHandler) *DataCleanStatDao {
	return &DataCleanStatDao{
		group:    "default",
		table:    "hg_data_clean_stat",
		columns:  dataCleanStatColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataCleanStatDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataCleanStatDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataCleanStatDao) Columns() DataCleanStatColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataCleanStatDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataCleanStatDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DataCleanStatDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
