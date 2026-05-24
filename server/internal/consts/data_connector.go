// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

func init() {
	dict.RegisterEnums("DataConnectorDirectionOptions", "数据源方向选项", DataConnectorDirectionOptions)
	dict.RegisterEnums("DataConnectorTypeOptions", "数据源类型选项", DataConnectorTypeOptions)
}

const (
	DataConnectorDirectionSource = "source" // 输入源
	DataConnectorDirectionSink   = "sink"   // 输出目标
)

var DataConnectorDirectionOptions = []*model.Option{
	dict.GenPrimaryOption(DataConnectorDirectionSource, "输入源"),
	dict.GenSuccessOption(DataConnectorDirectionSink, "输出目标"),
}

const (
	DataConnectorTypeHTTP          = "http"
	DataConnectorTypeKafka         = "kafka"
	DataConnectorTypeS3            = "s3"
	DataConnectorTypeLog           = "log"
	DataConnectorTypeManual        = "manual"
	DataConnectorTypeFeishu        = "feishu"
	DataConnectorTypeLark          = "lark"
	DataConnectorTypeElasticsearch = "elasticsearch"
	DataConnectorTypeOpensearch    = "opensearch"
	DataConnectorTypeSplunk        = "splunk"
)

var DataConnectorTypeOptions = []*model.Option{
	dict.GenInfoOption(DataConnectorTypeHTTP, "HTTP"),
	dict.GenPrimaryOption(DataConnectorTypeKafka, "Kafka"),
	dict.GenWarningOption(DataConnectorTypeS3, "S3"),
	dict.GenDefaultOption(DataConnectorTypeLog, "日志文件"),
	dict.GenDefaultOption(DataConnectorTypeManual, "手动录入"),
	dict.GenHashOption(DataConnectorTypeFeishu, "飞书"),
	dict.GenHashOption(DataConnectorTypeLark, "Lark"),
	dict.GenHashOption(DataConnectorTypeElasticsearch, "Elasticsearch"),
	dict.GenHashOption(DataConnectorTypeOpensearch, "OpenSearch"),
	dict.GenHashOption(DataConnectorTypeSplunk, "Splunk"),
}
