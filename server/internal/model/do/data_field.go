// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataField is the golang structure of table hg_data_field for DAO operations like Where/Data.
type DataField struct {
	g.Meta      `orm:"table:hg_data_field, do:true"`
	Id          any         // 字段ID
	TaskId      any         // 清洗任务ID
	ConnectorId any         // 输入连接ID
	FieldPath   any         // 字段路径
	FieldName   any         // 字段名称
	FieldType   any         // 字段类型
	FieldTypes  *gjson.Json // 出现过的字段类型集合
	SampleValue *gjson.Json // 样例值
	Status      any         // 状态
	Remark      any         // 备注
	FirstSeenAt *gtime.Time // 首次发现时间
	LastSeenAt  *gtime.Time // 最后发现时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 修改时间
}
