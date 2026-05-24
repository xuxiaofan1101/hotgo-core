// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataCleanTask is the golang structure for table data_clean_task.
type DataCleanTask struct {
	Id                 int64       `json:"id"                 orm:"id"                   description:"清洗任务ID"`
	SourceId           int64       `json:"sourceId"           orm:"source_id"            description:"输入连接ID"`
	TemplateId         int64       `json:"templateId"         orm:"template_id"          description:"标准字段模板ID"`
	Name               string      `json:"name"               orm:"name"                 description:"任务名称"`
	Code               string      `json:"code"               orm:"code"                 description:"任务编码"`
	EventType          string      `json:"eventType"          orm:"event_type"           description:"事件类型"`
	SourceConfig       *gjson.Json `json:"sourceConfig"       orm:"source_config"        description:"输入配置"`
	CleanEnabled       int         `json:"cleanEnabled"       orm:"clean_enabled"        description:"是否启用数据清洗"`
	CleanConfig        *gjson.Json `json:"cleanConfig"        orm:"clean_config"         description:"数据清洗配置"`
	CleanConfigVersion int         `json:"cleanConfigVersion" orm:"clean_config_version" description:"清洗配置版本"`
	FieldSchemaVersion int         `json:"fieldSchemaVersion" orm:"field_schema_version" description:"字段集合版本"`
	UnknownFieldPolicy string      `json:"unknownFieldPolicy" orm:"unknown_field_policy" description:"未知字段策略：selected_only strict passthrough drop"`
	CleanErrorPolicy   string      `json:"cleanErrorPolicy"   orm:"clean_error_policy"   description:"清洗失败策略：skip keep_raw stop_task"`
	SinkConfig         *gjson.Json `json:"sinkConfig"         orm:"sink_config"          description:"输出配置"`
	Status             int         `json:"status"             orm:"status"               description:"状态"`
	Remark             string      `json:"remark"             orm:"remark"               description:"备注"`
	CreatedBy          int64       `json:"createdBy"          orm:"created_by"           description:"创建者"`
	UpdatedBy          int64       `json:"updatedBy"          orm:"updated_by"           description:"更新者"`
	DeletedBy          int64       `json:"deletedBy"          orm:"deleted_by"           description:"删除者"`
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:"创建时间"`
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:"修改时间"`
	DeletedAt          *gtime.Time `json:"deletedAt"          orm:"deleted_at"           description:"删除时间"`
}
