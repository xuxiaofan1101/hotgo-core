// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataCleanTaskDao is the data access object for the table hg_data_clean_task.
type DataCleanTaskDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  DataCleanTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// DataCleanTaskColumns defines and stores column names for the table hg_data_clean_task.
type DataCleanTaskColumns struct {
	Id                 string // 清洗任务ID
	SourceId           string // 输入连接ID
	TemplateId         string // 标准字段模板ID
	Name               string // 任务名称
	Code               string // 任务编码
	EventType          string // 事件类型
	SourceConfig       string // 输入配置
	CleanEnabled       string // 是否启用数据清洗
	CleanConfig        string // 数据清洗配置
	CleanConfigVersion string // 清洗配置版本
	FieldSchemaVersion string // 字段集合版本
	UnknownFieldPolicy string // 未知字段策略：selected_only strict passthrough drop
	CleanErrorPolicy   string // 清洗失败策略：skip keep_raw stop_task
	SinkConfig         string // 输出配置
	Status             string // 状态
	Remark             string // 备注
	CreatedBy          string // 创建者
	UpdatedBy          string // 更新者
	DeletedBy          string // 删除者
	CreatedAt          string // 创建时间
	UpdatedAt          string // 修改时间
	DeletedAt          string // 删除时间
}

// dataCleanTaskColumns holds the columns for the table hg_data_clean_task.
var dataCleanTaskColumns = DataCleanTaskColumns{
	Id:                 "id",
	SourceId:           "source_id",
	TemplateId:         "template_id",
	Name:               "name",
	Code:               "code",
	EventType:          "event_type",
	SourceConfig:       "source_config",
	CleanEnabled:       "clean_enabled",
	CleanConfig:        "clean_config",
	CleanConfigVersion: "clean_config_version",
	FieldSchemaVersion: "field_schema_version",
	UnknownFieldPolicy: "unknown_field_policy",
	CleanErrorPolicy:   "clean_error_policy",
	SinkConfig:         "sink_config",
	Status:             "status",
	Remark:             "remark",
	CreatedBy:          "created_by",
	UpdatedBy:          "updated_by",
	DeletedBy:          "deleted_by",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewDataCleanTaskDao creates and returns a new DAO object for table data access.
func NewDataCleanTaskDao(handlers ...gdb.ModelHandler) *DataCleanTaskDao {
	return &DataCleanTaskDao{
		group:    "default",
		table:    "hg_data_clean_task",
		columns:  dataCleanTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataCleanTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataCleanTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataCleanTaskDao) Columns() DataCleanTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataCleanTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataCleanTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DataCleanTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
