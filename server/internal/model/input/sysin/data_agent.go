package sysin

import (
	"context"

	"hotgo/internal/consts"
	"hotgo/internal/model/input/form"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type DataAgentListInp struct {
	form.PageReq
	Keyword        string        `json:"keyword"        dc:"关键词"`
	RegisterStatus string        `json:"registerStatus" dc:"注册状态"`
	DispatchStatus string        `json:"dispatchStatus" dc:"调度状态"`
	OnlineStatus   string        `json:"onlineStatus"   dc:"在线状态"`
	CreatedAt      []*gtime.Time `json:"createdAt"      dc:"创建时间"`
}

func (in *DataAgentListInp) Filter(ctx context.Context) (err error) {
	return nil
}

type DataAgentListModel struct {
	Id              int64       `json:"id"              dc:"节点ID"`
	AgentId         string      `json:"agentId"         dc:"Agent唯一标识"`
	Name            string      `json:"name"            dc:"节点名称"`
	Hostname        string      `json:"hostname"        dc:"主机名"`
	Version         string      `json:"version"         dc:"Agent版本"`
	AgentIp         []string    `json:"agentIp"         dc:"Agent上报IP列表"`
	RegisterStatus  string      `json:"registerStatus"  dc:"注册状态"`
	DispatchStatus  string      `json:"dispatchStatus"  dc:"调度状态"`
	DispatchAllowed bool        `json:"dispatchAllowed" dc:"是否允许调度"`
	OnlineStatus    string      `json:"onlineStatus"    dc:"在线状态"`
	Online          bool        `json:"online"          dc:"是否在线"`
	LastSeenAt      *gtime.Time `json:"lastSeenAt"      dc:"最近心跳时间"`
	ApprovedBy      int64       `json:"approvedBy"      dc:"批准人"`
	ApprovedAt      *gtime.Time `json:"approvedAt"      dc:"批准时间"`
	DisabledAt      *gtime.Time `json:"disabledAt"      dc:"禁用/拒绝/吊销时间"`
	Remark          string      `json:"remark"          dc:"备注"`
	CreatedAt       *gtime.Time `json:"createdAt"       dc:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       dc:"修改时间"`
}

type DataAgentApproveInp struct {
	AgentId string `json:"agentId" v:"required#Agent ID不能为空" dc:"Agent唯一标识"`
	Remark  string `json:"remark"  dc:"备注"`
}

func (in *DataAgentApproveInp) Filter(ctx context.Context) (err error) {
	return nil
}

type DataAgentApproveModel struct{}

type DataAgentDispatchInp struct {
	AgentId        string `json:"agentId"        v:"required#Agent ID不能为空" dc:"Agent唯一标识"`
	DispatchStatus string `json:"dispatchStatus" dc:"调度状态"`
	Remark         string `json:"remark"         dc:"备注"`
}

func (in *DataAgentDispatchInp) Filter(ctx context.Context) (err error) {
	if in.DispatchStatus == "" {
		return gerror.New("调度状态不能为空")
	}
	if !validate.InSlice([]string{consts.DataAgentDispatchStatusEnabled, consts.DataAgentDispatchStatusDisabled}, in.DispatchStatus) {
		return gerror.New("调度状态不正确")
	}
	return nil
}

type DataAgentDispatchModel struct{}

type DataAgentRejectInp struct {
	AgentId        string `json:"agentId"        v:"required#Agent ID不能为空" dc:"Agent唯一标识"`
	RegisterStatus string `json:"registerStatus" dc:"注册状态"`
	Remark         string `json:"remark"         dc:"备注"`
}

func (in *DataAgentRejectInp) Filter(ctx context.Context) (err error) {
	if in.RegisterStatus == "" {
		return gerror.New("注册状态不能为空")
	}
	if !validate.InSlice([]string{consts.DataAgentRegisterStatusRejected, consts.DataAgentRegisterStatusRevoked}, in.RegisterStatus) {
		return gerror.New("注册状态不正确")
	}
	return nil
}

type DataAgentRejectModel struct{}
