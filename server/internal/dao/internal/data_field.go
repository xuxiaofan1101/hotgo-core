// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataFieldDao is the data access object for the table hg_data_field.
type DataFieldDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DataFieldColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DataFieldColumns defines and stores column names for the table hg_data_field.
type DataFieldColumns struct {
	Id          string // 字段ID
	TaskId      string // 清洗任务ID
	ConnectorId string // 输入连接ID
	FieldPath   string // 字段路径
	FieldName   string // 字段名称
	FieldType   string // 字段类型
	FieldTypes  string // 出现过的字段类型集合
	SampleValue string // 样例值
	Status      string // 状态
	Remark      string // 备注
	FirstSeenAt string // 首次发现时间
	LastSeenAt  string // 最后发现时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 修改时间
}

// dataFieldColumns holds the columns for the table hg_data_field.
var dataFieldColumns = DataFieldColumns{
	Id:          "id",
	TaskId:      "task_id",
	ConnectorId: "connector_id",
	FieldPath:   "field_path",
	FieldName:   "field_name",
	FieldType:   "field_type",
	FieldTypes:  "field_types",
	SampleValue: "sample_value",
	Status:      "status",
	Remark:      "remark",
	FirstSeenAt: "first_seen_at",
	LastSeenAt:  "last_seen_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewDataFieldDao creates and returns a new DAO object for table data access.
func NewDataFieldDao(handlers ...gdb.ModelHandler) *DataFieldDao {
	return &DataFieldDao{
		group:    "default",
		table:    "hg_data_field",
		columns:  dataFieldColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataFieldDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataFieldDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataFieldDao) Columns() DataFieldColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataFieldDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataFieldDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DataFieldDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
