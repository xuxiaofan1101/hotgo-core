// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataFieldTemplate is the golang structure of table hg_data_field_template for DAO operations like Where/Data.
type DataFieldTemplate struct {
	g.Meta    `orm:"table:hg_data_field_template, do:true"`
	Id        any         // 字段模板ID
	Name      any         // 模板名称
	Code      any         // 模板编码
	EventType any         // 事件类型
	Fields    *gjson.Json // 标准字段定义
	Examples  *gjson.Json // 标准样例数据
	Status    any         // 状态
	Remark    any         // 备注
	CreatedBy any         // 创建者
	UpdatedBy any         // 更新者
	DeletedBy any         // 删除者
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 修改时间
	DeletedAt *gtime.Time // 删除时间
}
