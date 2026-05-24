// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataCleanStat is the golang structure for table data_clean_stat.
type DataCleanStat struct {
	Id               int64       `json:"id"               orm:"id"                 description:"清洗统计ID"`
	TaskId           int64       `json:"taskId"           orm:"task_id"            description:"清洗任务ID"`
	ConnectorId      int64       `json:"connectorId"      orm:"connector_id"       description:"输入连接ID"`
	EventType        string      `json:"eventType"        orm:"event_type"         description:"事件类型"`
	WindowStart      *gtime.Time `json:"windowStart"      orm:"window_start"       description:"统计窗口开始时间"`
	WindowEnd        *gtime.Time `json:"windowEnd"        orm:"window_end"         description:"统计窗口结束时间"`
	TotalCount       int64       `json:"totalCount"       orm:"total_count"        description:"总处理数"`
	SuccessCount     int64       `json:"successCount"     orm:"success_count"      description:"成功数"`
	FailedCount      int64       `json:"failedCount"      orm:"failed_count"       description:"失败数"`
	CleanFailedCount int64       `json:"cleanFailedCount" orm:"clean_failed_count" description:"清洗失败数"`
	DroppedCount     int64       `json:"droppedCount"     orm:"dropped_count"      description:"丢弃数"`
	LastError        string      `json:"lastError"        orm:"last_error"         description:"最近错误"`
	ErrorSamples     *gjson.Json `json:"errorSamples"     orm:"error_samples"      description:"错误样例"`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:"创建时间"`
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:"修改时间"`
}
