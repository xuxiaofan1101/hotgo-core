// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataFieldTemplate is the golang structure for table data_field_template.
type DataFieldTemplate struct {
	Id        int64       `json:"id"        orm:"id"         description:"字段模板ID"`
	Name      string      `json:"name"      orm:"name"       description:"模板名称"`
	Code      string      `json:"code"      orm:"code"       description:"模板编码"`
	EventType string      `json:"eventType" orm:"event_type" description:"事件类型"`
	Fields    *gjson.Json `json:"fields"    orm:"fields"     description:"标准字段定义"`
	Examples  *gjson.Json `json:"examples"  orm:"examples"   description:"标准样例数据"`
	Status    int         `json:"status"    orm:"status"     description:"状态"`
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`
	CreatedBy int64       `json:"createdBy" orm:"created_by" description:"创建者"`
	UpdatedBy int64       `json:"updatedBy" orm:"updated_by" description:"更新者"`
	DeletedBy int64       `json:"deletedBy" orm:"deleted_by" description:"删除者"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"修改时间"`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`
}
