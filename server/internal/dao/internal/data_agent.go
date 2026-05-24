// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DataAgentDao is the data access object for the table hg_data_agent.
type DataAgentDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DataAgentColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DataAgentColumns defines and stores column names for the table hg_data_agent.
type DataAgentColumns struct {
	Id             string // 节点ID
	AgentId        string // Agent唯一标识
	Name           string // 节点名称
	Hostname       string // 主机名
	Version        string // Agent版本
	AgentIp        string // Agent上报IP列表
	RegisterStatus string // 注册状态：pending approved rejected revoked
	DispatchStatus string // 调度状态：enabled disabled
	LastSeenAt     string // 最近心跳时间
	ApprovedBy     string // 批准人
	ApprovedAt     string // 批准时间
	DisabledAt     string // 禁用/拒绝/吊销时间
	Remark         string // 备注
	CreatedAt      string // 创建时间
	UpdatedAt      string // 修改时间
}

// dataAgentColumns holds the columns for the table hg_data_agent.
var dataAgentColumns = DataAgentColumns{
	Id:             "id",
	AgentId:        "agent_id",
	Name:           "name",
	Hostname:       "hostname",
	Version:        "version",
	AgentIp:        "agent_ip",
	RegisterStatus: "register_status",
	DispatchStatus: "dispatch_status",
	LastSeenAt:     "last_seen_at",
	ApprovedBy:     "approved_by",
	ApprovedAt:     "approved_at",
	DisabledAt:     "disabled_at",
	Remark:         "remark",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewDataAgentDao creates and returns a new DAO object for table data access.
func NewDataAgentDao(handlers ...gdb.ModelHandler) *DataAgentDao {
	return &DataAgentDao{
		group:    "default",
		table:    "hg_data_agent",
		columns:  dataAgentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DataAgentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DataAgentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DataAgentDao) Columns() DataAgentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DataAgentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DataAgentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DataAgentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
