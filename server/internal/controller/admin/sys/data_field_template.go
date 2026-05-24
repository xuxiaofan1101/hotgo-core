package sys

import (
	"context"

	"hotgo/api/admin/datafieldtemplate"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var DataFieldTemplate = cDataFieldTemplate{}

type cDataFieldTemplate struct{}

func (c *cDataFieldTemplate) List(ctx context.Context, req *datafieldtemplate.ListReq) (res *datafieldtemplate.ListRes, err error) {
	list, totalCount, err := service.SysDataFieldTemplate().List(ctx, &req.DataFieldTemplateListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.DataFieldTemplateListModel{}
	}
	res = new(datafieldtemplate.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

func (c *cDataFieldTemplate) Edit(ctx context.Context, req *datafieldtemplate.EditReq) (res *datafieldtemplate.EditRes, err error) {
	err = service.SysDataFieldTemplate().Edit(ctx, &req.DataFieldTemplateEditInp)
	return
}

func (c *cDataFieldTemplate) View(ctx context.Context, req *datafieldtemplate.ViewReq) (res *datafieldtemplate.ViewRes, err error) {
	data, err := service.SysDataFieldTemplate().View(ctx, &req.DataFieldTemplateViewInp)
	if err != nil {
		return
	}
	res = new(datafieldtemplate.ViewRes)
	res.DataFieldTemplateViewModel = data
	return
}

func (c *cDataFieldTemplate) Delete(ctx context.Context, req *datafieldtemplate.DeleteReq) (res *datafieldtemplate.DeleteRes, err error) {
	err = service.SysDataFieldTemplate().Delete(ctx, &req.DataFieldTemplateDeleteInp)
	return
}

func (c *cDataFieldTemplate) Status(ctx context.Context, req *datafieldtemplate.StatusReq) (res *datafieldtemplate.StatusRes, err error) {
	err = service.SysDataFieldTemplate().Status(ctx, &req.DataFieldTemplateStatusInp)
	return
}
