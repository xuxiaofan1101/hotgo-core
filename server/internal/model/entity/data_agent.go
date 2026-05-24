// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataAgent is the golang structure for table data_agent.
type DataAgent struct {
	Id             int64       `json:"id"             orm:"id"              description:"节点ID"`
	AgentId        string      `json:"agentId"        orm:"agent_id"        description:"Agent唯一标识"`
	Name           string      `json:"name"           orm:"name"            description:"节点名称"`
	Hostname       string      `json:"hostname"       orm:"hostname"        description:"主机名"`
	Version        string      `json:"version"        orm:"version"         description:"Agent版本"`
	AgentIp        *gjson.Json `json:"agentIp"        orm:"agent_ip"        description:"Agent上报IP列表"`
	RegisterStatus string      `json:"registerStatus" orm:"register_status" description:"注册状态：pending approved rejected revoked"`
	DispatchStatus string      `json:"dispatchStatus" orm:"dispatch_status" description:"调度状态：enabled disabled"`
	LastSeenAt     *gtime.Time `json:"lastSeenAt"     orm:"last_seen_at"    description:"最近心跳时间"`
	ApprovedBy     int64       `json:"approvedBy"     orm:"approved_by"     description:"批准人"`
	ApprovedAt     *gtime.Time `json:"approvedAt"     orm:"approved_at"     description:"批准时间"`
	DisabledAt     *gtime.Time `json:"disabledAt"     orm:"disabled_at"     description:"禁用/拒绝/吊销时间"`
	Remark         string      `json:"remark"         orm:"remark"          description:"备注"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"修改时间"`
}
