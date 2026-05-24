// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataCleanStat is the golang structure of table hg_data_clean_stat for DAO operations like Where/Data.
type DataCleanStat struct {
	g.Meta           `orm:"table:hg_data_clean_stat, do:true"`
	Id               any         // 清洗统计ID
	TaskId           any         // 清洗任务ID
	ConnectorId      any         // 输入连接ID
	EventType        any         // 事件类型
	WindowStart      *gtime.Time // 统计窗口开始时间
	WindowEnd        *gtime.Time // 统计窗口结束时间
	TotalCount       any         // 总处理数
	SuccessCount     any         // 成功数
	FailedCount      any         // 失败数
	CleanFailedCount any         // 清洗失败数
	DroppedCount     any         // 丢弃数
	LastError        any         // 最近错误
	ErrorSamples     *gjson.Json // 错误样例
	CreatedAt        *gtime.Time // 创建时间
	UpdatedAt        *gtime.Time // 修改时间
}
