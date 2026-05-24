package sysin

import (
	"context"
	"encoding/json"
	"time"

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

const (
	DataAgentProtocolVersion = "v1"

	DataAgentMessageHello        = "agent.hello"
	DataAgentMessageHeartbeat    = "agent.heartbeat"
	DataAgentMessageFieldsReport = "agent.fields.report"
	DataAgentMessageStatsReport  = "agent.stats.report"
	DataAgentMessageError        = "agent.error"

	DataAgentMessageHelloAck     = "server.hello_ack"
	DataAgentMessageHeartbeatAck = "server.heartbeat_ack"
	DataAgentMessageFieldsAck    = "server.fields_ack"
	DataAgentMessageAck          = "server.ack"
	DataAgentMessageServerError  = "server.error"
)

type DataAgentEnvelope struct {
	Type      string          `json:"type"`
	RequestId string          `json:"requestId,omitempty"`
	SentAt    time.Time       `json:"sentAt,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type DataAgentSlotStatus struct {
	Total     int `json:"total"`
	Available int `json:"available"`
}

type DataAgentRuntimeMetrics struct {
	RunningShards int     `json:"runningShards"`
	ReadTps       float64 `json:"readTps"`
	CleanTps      float64 `json:"cleanTps"`
	OutputTps     float64 `json:"outputTps"`
}

type DataAgentHelloPayload struct {
	AgentId      string              `json:"agentId"`
	Hostname     string              `json:"hostname"`
	Capabilities []string            `json:"capabilities"`
	AgentIp      []string            `json:"agentIp"`
	Slots        DataAgentSlotStatus `json:"slots"`
	Version      string              `json:"version"`
}

type DataAgentHeartbeatPayload struct {
	AgentId string                  `json:"agentId"`
	AgentIp []string                `json:"agentIp"`
	Slots   DataAgentSlotStatus     `json:"slots"`
	Metrics DataAgentRuntimeMetrics `json:"metrics"`
	Time    time.Time               `json:"time"`
}

type DataAgentFieldSample struct {
	Path      string `json:"path"`
	Type      string `json:"type"`
	Sample    string `json:"sample"`
	Hash      string `json:"hash"`
	Count     int64  `json:"count"`
	NullCount int64  `json:"nullCount"`
}

type DataAgentFieldsReportPayload struct {
	TenantId string                 `json:"tenantId"`
	TaskId   int64                  `json:"taskId"`
	AgentId  string                 `json:"agentId"`
	Fields   []DataAgentFieldSample `json:"fields"`
}

type DataAgentStatsReportPayload struct {
	AgentId     string                 `json:"agentId"`
	WindowStart time.Time              `json:"windowStart"`
	WindowEnd   time.Time              `json:"windowEnd"`
	Metrics     []DataAgentStatsMetric `json:"metrics"`
}

type DataAgentStatsMetric struct {
	TenantId string `json:"tenantId"`
	TaskId   int64  `json:"taskId"`
	Name     string `json:"name"`
	GroupKey string `json:"groupKey,omitempty"`
	Value    int64  `json:"value"`
}
