// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataField is the golang structure for table data_field.
type DataField struct {
	Id          int64       `json:"id"          orm:"id"            description:"字段ID"`
	TaskId      int64       `json:"taskId"      orm:"task_id"       description:"清洗任务ID"`
	ConnectorId int64       `json:"connectorId" orm:"connector_id"  description:"输入连接ID"`
	FieldPath   string      `json:"fieldPath"   orm:"field_path"    description:"字段路径"`
	FieldName   string      `json:"fieldName"   orm:"field_name"    description:"字段名称"`
	FieldType   string      `json:"fieldType"   orm:"field_type"    description:"字段类型"`
	FieldTypes  *gjson.Json `json:"fieldTypes"  orm:"field_types"   description:"出现过的字段类型集合"`
	SampleValue *gjson.Json `json:"sampleValue" orm:"sample_value"  description:"样例值"`
	Status      int         `json:"status"      orm:"status"        description:"状态"`
	Remark      string      `json:"remark"      orm:"remark"        description:"备注"`
	FirstSeenAt *gtime.Time `json:"firstSeenAt" orm:"first_seen_at" description:"首次发现时间"`
	LastSeenAt  *gtime.Time `json:"lastSeenAt"  orm:"last_seen_at"  description:"最后发现时间"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:"修改时间"`
}
