// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataFieldTemplateDao is the data access object for the table hg_data_field_template.
type DataFieldTemplateDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  DataFieldTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// DataFieldTemplateColumns defines and stores column names for the table hg_data_field_template.
type DataFieldTemplateColumns struct {
	Id        string // 字段模板ID
	Name      string // 模板名称
	Code      string // 模板编码
	EventType string // 事件类型
	Fields    string // 标准字段定义
	Examples  string // 标准样例数据
	Status    string // 状态
	Remark    string // 备注
	CreatedBy string // 创建者
	UpdatedBy string // 更新者
	DeletedBy string // 删除者
	CreatedAt string // 创建时间
	UpdatedAt string // 修改时间
	DeletedAt string // 删除时间
}

// dataFieldTemplateColumns holds the columns for the table hg_data_field_template.
var dataFieldTemplateColumns = DataFieldTemplateColumns{
	Id:        "id",
	Name:      "name",
	Code:      "code",
	EventType: "event_type",
	Fields:    "fields",
	Examples:  "examples",
	Status:    "status",
	Remark:    "remark",
	CreatedBy: "created_by",
	UpdatedBy: "updated_by",
	DeletedBy: "deleted_by",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewDataFieldTemplateDao creates and returns a new DAO object for table data access.
func NewDataFieldTemplateDao(handlers ...gdb.ModelHandler) *DataFieldTemplateDao {
	return &DataFieldTemplateDao{
		group:    "default",
		table:    "hg_data_field_template",
		columns:  dataFieldTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataFieldTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataFieldTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataFieldTemplateDao) Columns() DataFieldTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataFieldTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataFieldTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DataFieldTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
