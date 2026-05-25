// Package dataconnector
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package dataconnector

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// ListReq 查询数据源列表
type ListReq struct {
	g.Meta `path:"/dataConnector/list" method:"get" tags:"数据源" summary:"获取数据源列表"`
	sysin.DataConnectorListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.DataConnectorListModel `json:"list" dc:"数据列表"`
}

// ViewReq 获取数据源指定信息
type ViewReq struct {
	g.Meta `path:"/dataConnector/view" method:"get" tags:"数据源" summary:"获取数据源指定信息"`
	sysin.DataConnectorViewInp
}

type ViewRes struct {
	*sysin.DataConnectorViewModel
}

// EditReq 修改/新增数据源
type EditReq struct {
	g.Meta `path:"/dataConnector/edit" method:"post" tags:"数据源" summary:"修改/新增数据源"`
	sysin.DataConnectorEditInp
}

type EditRes struct{}

// DeleteReq 删除数据源
type DeleteReq struct {
	g.Meta `path:"/dataConnector/delete" method:"post" tags:"数据源" summary:"删除数据源"`
	sysin.DataConnectorDeleteInp
}

type DeleteRes struct{}

// StatusReq 更新数据源状态
type StatusReq struct {
	g.Meta `path:"/dataConnector/status" method:"post" tags:"数据源" summary:"更新数据源状态"`
	sysin.DataConnectorStatusInp
}

type StatusRes struct{}

// TestReq 测试数据源配置
type TestReq struct {
	g.Meta `path:"/dataConnector/test" method:"post" tags:"数据源" summary:"测试数据源配置"`
	sysin.DataConnectorTestInp
}

type TestRes struct {
	*sysin.DataConnectorTestModel
}

// KafkaTopicsReq 获取Kafka Topic列表
type KafkaTopicsReq struct {
	g.Meta `path:"/dataConnector/kafkaTopics" method:"get" tags:"数据源" summary:"获取Kafka Topic列表"`
	sysin.DataConnectorKafkaTopicsInp
}

type KafkaTopicsRes struct {
	*sysin.DataConnectorKafkaTopicsModel
}
