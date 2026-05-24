// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/api/admin/dataconnector"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var (
	DataConnector = cDataConnector{}
)

type cDataConnector struct{}

// List 查看数据源列表
func (c *cDataConnector) List(ctx context.Context, req *dataconnector.ListReq) (res *dataconnector.ListRes, err error) {
	list, totalCount, err := service.SysDataConnector().List(ctx, &req.DataConnectorListInp)
	if err != nil {
		return
	}

	if list == nil {
		list = []*sysin.DataConnectorListModel{}
	}

	res = new(dataconnector.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Edit 更新数据源
func (c *cDataConnector) Edit(ctx context.Context, req *dataconnector.EditReq) (res *dataconnector.EditRes, err error) {
	err = service.SysDataConnector().Edit(ctx, &req.DataConnectorEditInp)
	return
}

// View 获取指定数据源信息
func (c *cDataConnector) View(ctx context.Context, req *dataconnector.ViewReq) (res *dataconnector.ViewRes, err error) {
	data, err := service.SysDataConnector().View(ctx, &req.DataConnectorViewInp)
	if err != nil {
		return
	}

	res = new(dataconnector.ViewRes)
	res.DataConnectorViewModel = data
	return
}

// Delete 删除数据源
func (c *cDataConnector) Delete(ctx context.Context, req *dataconnector.DeleteReq) (res *dataconnector.DeleteRes, err error) {
	err = service.SysDataConnector().Delete(ctx, &req.DataConnectorDeleteInp)
	return
}

// Status 更新数据源状态
func (c *cDataConnector) Status(ctx context.Context, req *dataconnector.StatusReq) (res *dataconnector.StatusRes, err error) {
	err = service.SysDataConnector().Status(ctx, &req.DataConnectorStatusInp)
	return
}

// Test 测试数据源配置
func (c *cDataConnector) Test(ctx context.Context, req *dataconnector.TestReq) (res *dataconnector.TestRes, err error) {
	data, err := service.SysDataConnector().Test(ctx, &req.DataConnectorTestInp)
	if err != nil {
		return
	}

	res = new(dataconnector.TestRes)
	res.DataConnectorTestModel = data
	return
}
