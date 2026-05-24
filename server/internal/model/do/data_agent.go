// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataAgent is the golang structure of table hg_data_agent for DAO operations like Where/Data.
type DataAgent struct {
	g.Meta         `orm:"table:hg_data_agent, do:true"`
	Id             any         // 节点ID
	AgentId        any         // Agent唯一标识
	Name           any         // 节点名称
	Hostname       any         // 主机名
	Version        any         // Agent版本
	AgentIp        *gjson.Json // Agent上报IP列表
	RegisterStatus any         // 注册状态：pending approved rejected revoked
	DispatchStatus any         // 调度状态：enabled disabled
	LastSeenAt     *gtime.Time // 最近心跳时间
	ApprovedBy     any         // 批准人
	ApprovedAt     *gtime.Time // 批准时间
	DisabledAt     *gtime.Time // 禁用/拒绝/吊销时间
	Remark         any         // 备注
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 修改时间
}
