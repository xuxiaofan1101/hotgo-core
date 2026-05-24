// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataConnectorDao is the data access object for the table hg_data_connector.
type DataConnectorDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  DataConnectorColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// DataConnectorColumns defines and stores column names for the table hg_data_connector.
type DataConnectorColumns struct {
	Id            string // 连接ID
	Name          string // 连接名称
	Code          string // 连接编码
	Direction     string // 连接方向：source=输入 sink=输出
	ConnectorType string // 连接类型：http,kafka,s3,log,manual,feishu,lark,elasticsearch,opensearch,splunk
	Config        string // 连接配置
	Status        string // 状态
	Remark        string // 备注
	CreatedBy     string // 创建者
	UpdatedBy     string // 更新者
	DeletedBy     string // 删除者
	CreatedAt     string // 创建时间
	UpdatedAt     string // 修改时间
	DeletedAt     string // 删除时间
}

// dataConnectorColumns holds the columns for the table hg_data_connector.
var dataConnectorColumns = DataConnectorColumns{
	Id:            "id",
	Name:          "name",
	Code:          "code",
	Direction:     "direction",
	ConnectorType: "connector_type",
	Config:        "config",
	Status:        "status",
	Remark:        "remark",
	CreatedBy:     "created_by",
	UpdatedBy:     "updated_by",
	DeletedBy:     "deleted_by",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewDataConnectorDao creates and returns a new DAO object for table data access.
func NewDataConnectorDao(handlers ...gdb.ModelHandler) *DataConnectorDao {
	return &DataConnectorDao{
		group:    "default",
		table:    "hg_data_connector",
		columns:  dataConnectorColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataConnectorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataConnectorDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataConnectorDao) Columns() DataConnectorColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataConnectorDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataConnectorDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DataConnectorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
