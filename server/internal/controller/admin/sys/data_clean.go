package sys

import (
	"context"

	"hotgo/api/admin/dataclean"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var DataClean = cDataClean{}

type cDataClean struct{}

func (c *cDataClean) List(ctx context.Context, req *dataclean.ListReq) (res *dataclean.ListRes, err error) {
	list, totalCount, err := service.SysDataClean().List(ctx, &req.DataCleanTaskListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.DataCleanTaskListModel{}
	}
	res = new(dataclean.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

func (c *cDataClean) Edit(ctx context.Context, req *dataclean.EditReq) (res *dataclean.EditRes, err error) {
	err = service.SysDataClean().Edit(ctx, &req.DataCleanTaskEditInp)
	return
}

func (c *cDataClean) View(ctx context.Context, req *dataclean.ViewReq) (res *dataclean.ViewRes, err error) {
	data, err := service.SysDataClean().View(ctx, &req.DataCleanTaskViewInp)
	if err != nil {
		return
	}
	res = new(dataclean.ViewRes)
	res.DataCleanTaskViewModel = data
	return
}

func (c *cDataClean) Delete(ctx context.Context, req *dataclean.DeleteReq) (res *dataclean.DeleteRes, err error) {
	err = service.SysDataClean().Delete(ctx, &req.DataCleanTaskDeleteInp)
	return
}

func (c *cDataClean) Status(ctx context.Context, req *dataclean.StatusReq) (res *dataclean.StatusRes, err error) {
	err = service.SysDataClean().Status(ctx, &req.DataCleanTaskStatusInp)
	return
}

func (c *cDataClean) Test(ctx context.Context, req *dataclean.TestReq) (res *dataclean.TestRes, err error) {
	data, err := service.SysDataClean().Test(ctx, &req.DataCleanTaskTestInp)
	if err != nil {
		return
	}
	res = new(dataclean.TestRes)
	res.DataCleanTaskTestModel = data
	return
}

func (c *cDataClean) Sample(ctx context.Context, req *dataclean.SampleReq) (res *dataclean.SampleRes, err error) {
	data, err := service.SysDataClean().Sample(ctx, &req.DataCleanTaskSampleInp)
	if err != nil {
		return
	}
	res = new(dataclean.SampleRes)
	res.DataCleanTaskSampleModel = data
	return
}

func (c *cDataClean) FieldList(ctx context.Context, req *dataclean.FieldListReq) (res *dataclean.FieldListRes, err error) {
	list, total, err := service.SysDataClean().FieldList(ctx, &req.DataFieldListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.DataFieldListModel{}
	}
	res = &dataclean.FieldListRes{List: list, Total: total}
	return
}
