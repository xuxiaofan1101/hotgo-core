package sys

import (
	"context"

	"hotgo/api/admin/dataagent"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var DataAgent = cDataAgent{}

type cDataAgent struct{}

func (c *cDataAgent) List(ctx context.Context, req *dataagent.ListReq) (res *dataagent.ListRes, err error) {
	list, totalCount, err := service.SysDataAgent().List(ctx, &req.DataAgentListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.DataAgentListModel{}
	}
	res = new(dataagent.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

func (c *cDataAgent) Approve(ctx context.Context, req *dataagent.ApproveReq) (res *dataagent.ApproveRes, err error) {
	err = service.SysDataAgent().Approve(ctx, &req.DataAgentApproveInp)
	return
}

func (c *cDataAgent) Dispatch(ctx context.Context, req *dataagent.DispatchReq) (res *dataagent.DispatchRes, err error) {
	err = service.SysDataAgent().Dispatch(ctx, &req.DataAgentDispatchInp)
	return
}

func (c *cDataAgent) Reject(ctx context.Context, req *dataagent.RejectReq) (res *dataagent.RejectRes, err error) {
	err = service.SysDataAgent().Reject(ctx, &req.DataAgentRejectInp)
	return
}
