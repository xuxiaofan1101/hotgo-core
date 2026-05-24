// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/library/dict"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataConnectorUpdateFields 修改数据源字段过滤
type DataConnectorUpdateFields struct {
	Name          string      `json:"name"          dc:"连接名称"`
	Code          string      `json:"code"          dc:"连接编码"`
	Direction     string      `json:"direction"     dc:"连接方向"`
	ConnectorType string      `json:"connectorType" dc:"连接类型"`
	Config        *gjson.Json `json:"config"        dc:"连接配置"`
	Status        int         `json:"status"        dc:"状态"`
	Remark        string      `json:"remark"        dc:"备注"`
	UpdatedBy     int64       `json:"updatedBy"     dc:"更新者"`
}

// DataConnectorInsertFields 新增数据源字段过滤
type DataConnectorInsertFields struct {
	Name          string      `json:"name"          dc:"连接名称"`
	Code          string      `json:"code"          dc:"连接编码"`
	Direction     string      `json:"direction"     dc:"连接方向"`
	ConnectorType string      `json:"connectorType" dc:"连接类型"`
	Config        *gjson.Json `json:"config"        dc:"连接配置"`
	Status        int         `json:"status"        dc:"状态"`
	Remark        string      `json:"remark"        dc:"备注"`
	CreatedBy     int64       `json:"createdBy"     dc:"创建者"`
}

// DataConnectorEditInp 修改/新增数据源
type DataConnectorEditInp struct {
	entity.DataConnector
}

func (in *DataConnectorEditInp) Filter(ctx context.Context) (err error) {
	if err := g.Validator().Rules("required").Data(in.Name).Messages("连接名称不能为空").Run(ctx); err != nil {
		return err.Current()
	}

	if err := g.Validator().Rules("required").Data(in.Code).Messages("连接编码不能为空").Run(ctx); err != nil {
		return err.Current()
	}

	if err := g.Validator().Rules("required").Data(in.Direction).Messages("连接方向不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	if !dict.HasOptionKey(consts.DataConnectorDirectionOptions, in.Direction) {
		return gerror.New("连接方向不正确")
	}

	if err := g.Validator().Rules("required").Data(in.ConnectorType).Messages("连接类型不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	if !dict.HasOptionKey(consts.DataConnectorTypeOptions, in.ConnectorType) {
		return gerror.New("连接类型不正确")
	}

	if in.Status <= 0 {
		return gerror.New("状态不能为空")
	}
	if !validate.InSlice(consts.StatusSlice, in.Status) {
		return gerror.New("状态不正确")
	}
	return
}

type DataConnectorEditModel struct{}

// DataConnectorDeleteInp 删除数据源
type DataConnectorDeleteInp struct {
	Id interface{} `json:"id" v:"required#连接ID不能为空" dc:"连接ID"`
}

func (in *DataConnectorDeleteInp) Filter(ctx context.Context) (err error) {
	return
}

type DataConnectorDeleteModel struct{}

// DataConnectorViewInp 获取指定数据源信息
type DataConnectorViewInp struct {
	Id int64 `json:"id" v:"required#连接ID不能为空" dc:"连接ID"`
}

func (in *DataConnectorViewInp) Filter(ctx context.Context) (err error) {
	return
}

type DataConnectorViewModel struct {
	entity.DataConnector
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma" dc:"创建者摘要信息"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma" dc:"更新者摘要信息"`
}

// DataConnectorListInp 获取数据源列表
type DataConnectorListInp struct {
	form.PageReq
	Id            int64         `json:"id"            dc:"连接ID"`
	Keyword       string        `json:"keyword"       dc:"关键词"`
	Name          string        `json:"name"          dc:"连接名称"`
	Code          string        `json:"code"          dc:"连接编码"`
	Direction     string        `json:"direction"     dc:"连接方向"`
	ConnectorType string        `json:"connectorType" dc:"连接类型"`
	Status        int           `json:"status"        dc:"状态"`
	CreatedBy     string        `json:"createdBy"     dc:"创建者"`
	CreatedAt     []*gtime.Time `json:"createdAt"     dc:"创建时间"`
}

func (in *DataConnectorListInp) Filter(ctx context.Context) (err error) {
	return
}

type DataConnectorListModel struct {
	Id             int64             `json:"id"            dc:"连接ID"`
	Name           string            `json:"name"          dc:"连接名称"`
	Code           string            `json:"code"          dc:"连接编码"`
	Direction      string            `json:"direction"     dc:"连接方向"`
	ConnectorType  string            `json:"connectorType" dc:"连接类型"`
	Status         int               `json:"status"        dc:"状态"`
	Remark         string            `json:"remark"        dc:"备注"`
	CreatedBy      int64             `json:"createdBy"     dc:"创建者"`
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma" dc:"创建者摘要信息"`
	UpdatedBy      int64             `json:"updatedBy"     dc:"更新者"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma" dc:"更新者摘要信息"`
	CreatedAt      *gtime.Time       `json:"createdAt"     dc:"创建时间"`
	UpdatedAt      *gtime.Time       `json:"updatedAt"     dc:"修改时间"`
}

// DataConnectorStatusInp 更新数据源状态
type DataConnectorStatusInp struct {
	Id     int64 `json:"id" v:"required#连接ID不能为空" dc:"连接ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *DataConnectorStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		return gerror.New("连接ID不能为空")
	}

	if in.Status <= 0 {
		return gerror.New("状态不能为空")
	}

	if !validate.InSlice(consts.StatusSlice, in.Status) {
		return gerror.New("状态不正确")
	}
	return
}

type DataConnectorStatusModel struct{}

// DataConnectorTestInp 测试数据源配置
type DataConnectorTestInp struct {
	Direction     string      `json:"direction"     dc:"连接方向"`
	ConnectorType string      `json:"connectorType" dc:"连接类型"`
	Config        *gjson.Json `json:"config"        dc:"连接配置"`
	Message       string      `json:"message"       dc:"测试消息"`
}

func (in *DataConnectorTestInp) Filter(ctx context.Context) (err error) {
	if err := g.Validator().Rules("required").Data(in.Direction).Messages("连接方向不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	if !dict.HasOptionKey(consts.DataConnectorDirectionOptions, in.Direction) {
		return gerror.New("连接方向不正确")
	}

	if err := g.Validator().Rules("required").Data(in.ConnectorType).Messages("连接类型不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	if !dict.HasOptionKey(consts.DataConnectorTypeOptions, in.ConnectorType) {
		return gerror.New("连接类型不正确")
	}
	return
}

type DataConnectorTestModel struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"测试结果"`
}
