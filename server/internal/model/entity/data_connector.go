// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataConnector is the golang structure for table data_connector.
type DataConnector struct {
	Id            int64       `json:"id"            orm:"id"             description:"连接ID"`
	Name          string      `json:"name"          orm:"name"           description:"连接名称"`
	Code          string      `json:"code"          orm:"code"           description:"连接编码"`
	Direction     string      `json:"direction"     orm:"direction"      description:"连接方向：source=输入 sink=输出"`
	ConnectorType string      `json:"connectorType" orm:"connector_type" description:"连接类型：http,kafka,s3,log,manual,feishu,lark,elasticsearch,opensearch,splunk"`
	Config        *gjson.Json `json:"config"        orm:"config"         description:"连接配置"`
	Status        int         `json:"status"        orm:"status"         description:"状态"`
	Remark        string      `json:"remark"        orm:"remark"         description:"备注"`
	CreatedBy     int64       `json:"createdBy"     orm:"created_by"     description:"创建者"`
	UpdatedBy     int64       `json:"updatedBy"     orm:"updated_by"     description:"更新者"`
	DeletedBy     int64       `json:"deletedBy"     orm:"deleted_by"     description:"删除者"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"修改时间"`
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     description:"删除时间"`
}
