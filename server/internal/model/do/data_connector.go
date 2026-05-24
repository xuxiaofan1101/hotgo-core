// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataConnector is the golang structure of table hg_data_connector for DAO operations like Where/Data.
type DataConnector struct {
	g.Meta        `orm:"table:hg_data_connector, do:true"`
	Id            any         // 连接ID
	Name          any         // 连接名称
	Code          any         // 连接编码
	Direction     any         // 连接方向：source=输入 sink=输出
	ConnectorType any         // 连接类型：http,kafka,s3,log,manual,feishu,lark,elasticsearch,opensearch,splunk
	Config        *gjson.Json // 连接配置
	Status        any         // 状态
	Remark        any         // 备注
	CreatedBy     any         // 创建者
	UpdatedBy     any         // 更新者
	DeletedBy     any         // 删除者
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 修改时间
	DeletedAt     *gtime.Time // 删除时间
}
