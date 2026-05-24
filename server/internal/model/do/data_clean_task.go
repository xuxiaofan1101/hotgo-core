// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataCleanTask is the golang structure of table hg_data_clean_task for DAO operations like Where/Data.
type DataCleanTask struct {
	g.Meta             `orm:"table:hg_data_clean_task, do:true"`
	Id                 any         // 清洗任务ID
	SourceId           any         // 输入连接ID
	TemplateId         any         // 标准字段模板ID
	Name               any         // 任务名称
	Code               any         // 任务编码
	EventType          any         // 事件类型
	SourceConfig       *gjson.Json // 输入配置
	CleanEnabled       any         // 是否启用数据清洗
	CleanConfig        *gjson.Json // 数据清洗配置
	CleanConfigVersion any         // 清洗配置版本
	FieldSchemaVersion any         // 字段集合版本
	UnknownFieldPolicy any         // 未知字段策略：selected_only strict passthrough drop
	CleanErrorPolicy   any         // 清洗失败策略：skip keep_raw stop_task
	SinkConfig         *gjson.Json // 输出配置
	Status             any         // 状态
	Remark             any         // 备注
	CreatedBy          any         // 创建者
	UpdatedBy          any         // 更新者
	DeletedBy          any         // 删除者
	CreatedAt          *gtime.Time // 创建时间
	UpdatedAt          *gtime.Time // 修改时间
	DeletedAt          *gtime.Time // 删除时间
}
