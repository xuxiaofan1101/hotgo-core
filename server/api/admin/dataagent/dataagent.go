package dataagent

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta `path:"/dataAgent/list" method:"get" tags:"节点管理" summary:"获取Agent节点列表"`
	sysin.DataAgentListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.DataAgentListModel `json:"list" dc:"数据列表"`
}

type ApproveReq struct {
	g.Meta `path:"/dataAgent/approve" method:"post" tags:"节点管理" summary:"批准Agent注册并允许调度"`
	sysin.DataAgentApproveInp
}

type ApproveRes struct{}

type DispatchReq struct {
	g.Meta `path:"/dataAgent/dispatch" method:"post" tags:"节点管理" summary:"设置Agent调度状态"`
	sysin.DataAgentDispatchInp
}

type DispatchRes struct{}

type RejectReq struct {
	g.Meta `path:"/dataAgent/reject" method:"post" tags:"节点管理" summary:"拒绝或吊销Agent"`
	sysin.DataAgentRejectInp
}

type RejectRes struct{}
