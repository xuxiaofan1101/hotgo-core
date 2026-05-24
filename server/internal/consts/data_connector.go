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
	dict.RegisterEnums("DataCleanErrorPolicyOptions", "数据清洗失败策略选项", DataCleanErrorPolicyOptions)
	dict.RegisterEnums("DataUnknownFieldPolicyOptions", "数据清洗未知字段策略选项", DataUnknownFieldPolicyOptions)
	dict.RegisterEnums("DataFieldTypeOptions", "数据字段类型选项", DataFieldTypeOptions)
	dict.RegisterEnums("DataCleanActionOptions", "数据清洗动作选项", DataCleanActionOptions)
	dict.RegisterEnums("DataCleanFilterOperatorOptions", "数据清洗过滤操作符选项", DataCleanFilterOperatorOptions)
	dict.RegisterEnums("DataSinkTargetTypeOptions", "数据清洗输出目标类型选项", DataSinkTargetTypeOptions)
	dict.RegisterEnums("DataAgentRegisterStatusOptions", "Agent注册状态选项", DataAgentRegisterStatusOptions)
	dict.RegisterEnums("DataAgentDispatchStatusOptions", "Agent调度状态选项", DataAgentDispatchStatusOptions)
	dict.RegisterEnums("DataAgentOnlineStatusOptions", "Agent在线状态选项", DataAgentOnlineStatusOptions)
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

const (
	DataCleanErrorPolicySkip     = "skip"
	DataCleanErrorPolicyKeepRaw  = "keep_raw"
	DataCleanErrorPolicyStopTask = "stop_task"
)

var DataCleanErrorPolicyOptions = []*model.Option{
	dict.GenWarningOption(DataCleanErrorPolicySkip, "跳过异常数据"),
	dict.GenInfoOption(DataCleanErrorPolicyKeepRaw, "保留原始数据"),
	dict.GenErrorOption(DataCleanErrorPolicyStopTask, "停止当前任务"),
}

const (
	DataUnknownFieldPolicySelectedOnly = "selected_only"
	DataUnknownFieldPolicyStrict       = "strict"
	DataUnknownFieldPolicyPassthrough  = "passthrough"
	DataUnknownFieldPolicyDrop         = "drop"
)

var DataUnknownFieldPolicyOptions = []*model.Option{
	dict.GenPrimaryOption(DataUnknownFieldPolicySelectedOnly, "仅输出已选择字段"),
	dict.GenWarningOption(DataUnknownFieldPolicyStrict, "发现未知字段时报错"),
	dict.GenInfoOption(DataUnknownFieldPolicyPassthrough, "未知字段透传"),
	dict.GenDefaultOption(DataUnknownFieldPolicyDrop, "未知字段丢弃"),
}

var DataFieldTypeOptions = []*model.Option{
	dict.GenInfoOption("string", "字符串"),
	dict.GenPrimaryOption("number", "数字"),
	dict.GenPrimaryOption("int", "整数"),
	dict.GenPrimaryOption("float", "小数"),
	dict.GenSuccessOption("boolean", "布尔"),
	dict.GenWarningOption("object", "对象"),
	dict.GenWarningOption("array", "数组"),
	dict.GenDefaultOption("null", "空值"),
}

var DataCleanActionOptions = []*model.Option{
	dict.GenInfoOption("trim", "去除空格"),
	dict.GenInfoOption("to_string", "转字符串"),
	dict.GenPrimaryOption("to_int", "转整数"),
	dict.GenPrimaryOption("to_float", "转小数"),
	dict.GenSuccessOption("to_bool", "转布尔"),
	dict.GenWarningOption("parse_json", "解析 JSON"),
	dict.GenWarningOption("mask", "字段脱敏"),
	dict.GenDefaultOption("default", "默认值"),
	dict.GenDefaultOption("set_field", "设置字段"),
	dict.GenDefaultOption("rename", "重命名字段"),
	dict.GenDefaultOption("regex_replace", "正则替换"),
	dict.GenErrorOption("remove", "删除字段"),
}

var DataCleanFilterOperatorOptions = []*model.Option{
	dict.GenDefaultOption("eq", "等于"),
	dict.GenDefaultOption("ne", "不等于"),
	dict.GenDefaultOption("contains", "包含"),
	dict.GenDefaultOption("not_contains", "不包含"),
	dict.GenDefaultOption("matches", "正则匹配"),
	dict.GenDefaultOption("empty", "为空"),
	dict.GenDefaultOption("not_empty", "不为空"),
	dict.GenDefaultOption("gt", "大于"),
	dict.GenDefaultOption("gte", "大于等于"),
	dict.GenDefaultOption("lt", "小于"),
	dict.GenDefaultOption("lte", "小于等于"),
	dict.GenDefaultOption("between", "区间"),
	dict.GenDefaultOption("in", "在列表中"),
	dict.GenDefaultOption("not_in", "不在列表中"),
}

const (
	DataSinkTargetTypeStrategyEngine = "strategy_engine"
	DataSinkTargetTypeSink           = "sink"
)

var DataSinkTargetTypeOptions = []*model.Option{
	dict.GenPrimaryOption(DataSinkTargetTypeStrategyEngine, "策略引擎"),
	dict.GenSuccessOption(DataSinkTargetTypeSink, "输出目标"),
}

const (
	DataAgentRegisterStatusPending  = "pending"
	DataAgentRegisterStatusApproved = "approved"
	DataAgentRegisterStatusRejected = "rejected"
	DataAgentRegisterStatusRevoked  = "revoked"
)

var DataAgentRegisterStatusOptions = []*model.Option{
	dict.GenWarningOption(DataAgentRegisterStatusPending, "待审批"),
	dict.GenSuccessOption(DataAgentRegisterStatusApproved, "已批准"),
	dict.GenErrorOption(DataAgentRegisterStatusRejected, "已拒绝"),
	dict.GenErrorOption(DataAgentRegisterStatusRevoked, "已吊销"),
}

const (
	DataAgentDispatchStatusEnabled  = "enabled"
	DataAgentDispatchStatusDisabled = "disabled"
)

var DataAgentDispatchStatusOptions = []*model.Option{
	dict.GenSuccessOption(DataAgentDispatchStatusEnabled, "已开启"),
	dict.GenWarningOption(DataAgentDispatchStatusDisabled, "已禁用"),
}

const (
	DataAgentOnlineStatusOnline  = "online"
	DataAgentOnlineStatusOffline = "offline"
)

var DataAgentOnlineStatusOptions = []*model.Option{
	dict.GenSuccessOption(DataAgentOnlineStatusOnline, "在线"),
	dict.GenDefaultOption(DataAgentOnlineStatusOffline, "离线"),
}
