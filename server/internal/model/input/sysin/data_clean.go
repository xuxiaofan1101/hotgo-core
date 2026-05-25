package sysin

import (
	"context"

	"hotgo/internal/consts"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DataCleanTaskUpdateFields 修改数据清洗任务字段过滤
type DataCleanTaskUpdateFields struct {
	SourceId           int64       `json:"sourceId"           dc:"输入连接ID"`
	TemplateId         int64       `json:"templateId"         dc:"标准字段模板ID"`
	Name               string      `json:"name"               dc:"任务名称"`
	Code               string      `json:"code"               dc:"任务编码"`
	EventType          string      `json:"eventType"          dc:"事件类型"`
	SourceConfig       *gjson.Json `json:"sourceConfig"       dc:"输入配置"`
	CleanEnabled       int         `json:"cleanEnabled"       dc:"是否启用数据清洗"`
	CleanConfig        *gjson.Json `json:"cleanConfig"        dc:"数据清洗配置"`
	CleanConfigVersion int         `json:"cleanConfigVersion" dc:"清洗配置版本"`
	UnknownFieldPolicy string      `json:"unknownFieldPolicy" dc:"未知字段策略"`
	CleanErrorPolicy   string      `json:"cleanErrorPolicy"   dc:"清洗失败策略"`
	SinkConfig         *gjson.Json `json:"sinkConfig"         dc:"输出配置"`
	Status             int         `json:"status"             dc:"状态"`
	Remark             string      `json:"remark"             dc:"备注"`
	UpdatedBy          int64       `json:"updatedBy"          dc:"更新者"`
}

// DataCleanTaskInsertFields 新增数据清洗任务字段过滤
type DataCleanTaskInsertFields struct {
	SourceId           int64       `json:"sourceId"           dc:"输入连接ID"`
	TemplateId         int64       `json:"templateId"         dc:"标准字段模板ID"`
	Name               string      `json:"name"               dc:"任务名称"`
	Code               string      `json:"code"               dc:"任务编码"`
	EventType          string      `json:"eventType"          dc:"事件类型"`
	SourceConfig       *gjson.Json `json:"sourceConfig"       dc:"输入配置"`
	CleanEnabled       int         `json:"cleanEnabled"       dc:"是否启用数据清洗"`
	CleanConfig        *gjson.Json `json:"cleanConfig"        dc:"数据清洗配置"`
	CleanConfigVersion int         `json:"cleanConfigVersion" dc:"清洗配置版本"`
	FieldSchemaVersion int         `json:"fieldSchemaVersion" dc:"字段集合版本"`
	UnknownFieldPolicy string      `json:"unknownFieldPolicy" dc:"未知字段策略"`
	CleanErrorPolicy   string      `json:"cleanErrorPolicy"   dc:"清洗失败策略"`
	SinkConfig         *gjson.Json `json:"sinkConfig"         dc:"输出配置"`
	Status             int         `json:"status"             dc:"状态"`
	Remark             string      `json:"remark"             dc:"备注"`
	CreatedBy          int64       `json:"createdBy"          dc:"创建者"`
}

// DataCleanTaskEditInp 修改/新增数据清洗任务
type DataCleanTaskEditInp struct {
	entity.DataCleanTask
}

func (in *DataCleanTaskEditInp) Filter(ctx context.Context) (err error) {
	if in.SourceId <= 0 {
		return gerror.New("请选择输入数据源")
	}
	if e := g.Validator().Rules("required").Data(in.Name).Messages("清洗任务名称不能为空").Run(ctx); e != nil {
		return e.Current()
	}
	if e := g.Validator().Rules("required").Data(in.Code).Messages("清洗任务编码不能为空").Run(ctx); e != nil {
		return e.Current()
	}
	if in.EventType == "" {
		in.EventType = "custom.event"
	}
	if in.UnknownFieldPolicy == "" {
		in.UnknownFieldPolicy = consts.DataUnknownFieldPolicySelectedOnly
	}
	if in.CleanErrorPolicy == "" {
		in.CleanErrorPolicy = consts.DataCleanErrorPolicySkip
	}
	if in.Status <= 0 {
		in.Status = consts.StatusDisable
	}
	if !validate.InSlice(consts.StatusSlice, in.Status) {
		return gerror.New("状态不正确")
	}
	if in.CleanEnabled != 0 && in.CleanEnabled != 1 {
		return gerror.New("清洗启用状态不正确")
	}
	return nil
}

type DataCleanTaskEditModel struct{}

type DataCleanTaskDeleteInp struct {
	Id interface{} `json:"id" v:"required#清洗任务ID不能为空" dc:"清洗任务ID"`
}

func (in *DataCleanTaskDeleteInp) Filter(ctx context.Context) (err error) { return }

type DataCleanTaskDeleteModel struct{}

type DataCleanTaskViewInp struct {
	Id int64 `json:"id" v:"required#清洗任务ID不能为空" dc:"清洗任务ID"`
}

func (in *DataCleanTaskViewInp) Filter(ctx context.Context) (err error) { return }

type DataCleanTaskViewModel struct {
	entity.DataCleanTask
	SourceName     string            `json:"sourceName"      dc:"输入源名称"`
	SourceCode     string            `json:"sourceCode"      dc:"输入源编码"`
	SourceType     string            `json:"sourceType"      dc:"输入源类型"`
	TemplateName   string            `json:"templateName"    dc:"模板名称"`
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma"  dc:"创建者摘要信息"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma"  dc:"更新者摘要信息"`
}

type DataCleanTaskListInp struct {
	form.PageReq
	Id         int64         `json:"id"            dc:"清洗任务ID"`
	Keyword    string        `json:"keyword"       dc:"关键词"`
	SourceId   int64         `json:"sourceId"      dc:"输入连接ID"`
	SourceType string        `json:"sourceType"    dc:"输入连接类型"`
	TemplateId int64         `json:"templateId"    dc:"标准字段模板ID"`
	EventType  string        `json:"eventType"     dc:"事件类型"`
	Topic      string        `json:"topic"         dc:"Topic"`
	GroupId    string        `json:"groupId"       dc:"Group ID"`
	Bucket     string        `json:"bucket"        dc:"Bucket"`
	Prefix     string        `json:"prefix"        dc:"Prefix"`
	Path       string        `json:"path"          dc:"日志路径"`
	Status     int           `json:"status"        dc:"状态"`
	CreatedBy  string        `json:"createdBy"     dc:"创建者"`
	CreatedAt  []*gtime.Time `json:"createdAt"     dc:"创建时间"`
}

func (in *DataCleanTaskListInp) Filter(ctx context.Context) (err error) { return }

type DataCleanTaskListModel struct {
	Id                 int64             `json:"id"                 dc:"清洗任务ID"`
	SourceId           int64             `json:"sourceId"           dc:"输入连接ID"`
	SourceName         string            `json:"sourceName"         dc:"输入源名称"`
	SourceCode         string            `json:"sourceCode"         dc:"输入源编码"`
	SourceType         string            `json:"sourceType"         dc:"输入源类型"`
	TemplateId         int64             `json:"templateId"         dc:"标准字段模板ID"`
	TemplateName       string            `json:"templateName"       dc:"模板名称"`
	Name               string            `json:"name"               dc:"任务名称"`
	Code               string            `json:"code"               dc:"任务编码"`
	EventType          string            `json:"eventType"          dc:"事件类型"`
	SourceConfig       *gjson.Json       `json:"sourceConfig"       dc:"输入配置"`
	CleanEnabled       int               `json:"cleanEnabled"       dc:"是否启用数据清洗"`
	CleanConfig        *gjson.Json       `json:"cleanConfig"        dc:"数据清洗配置"`
	CleanConfigVersion int               `json:"cleanConfigVersion" dc:"清洗配置版本"`
	FieldSchemaVersion int               `json:"fieldSchemaVersion" dc:"字段集合版本"`
	UnknownFieldPolicy string            `json:"unknownFieldPolicy" dc:"未知字段策略"`
	CleanErrorPolicy   string            `json:"cleanErrorPolicy"   dc:"清洗失败策略"`
	SinkConfig         *gjson.Json       `json:"sinkConfig"         dc:"输出配置"`
	Status             int               `json:"status"             dc:"状态"`
	Remark             string            `json:"remark"             dc:"备注"`
	CreatedBy          int64             `json:"createdBy"          dc:"创建者"`
	CreatedBySumma     *hook.MemberSumma `json:"createdBySumma"     dc:"创建者摘要信息"`
	UpdatedBy          int64             `json:"updatedBy"          dc:"更新者"`
	UpdatedBySumma     *hook.MemberSumma `json:"updatedBySumma"     dc:"更新者摘要信息"`
	CreatedAt          *gtime.Time       `json:"createdAt"          dc:"创建时间"`
	UpdatedAt          *gtime.Time       `json:"updatedAt"          dc:"修改时间"`
}

type DataCleanTaskStatusInp struct {
	Id     int64 `json:"id" v:"required#清洗任务ID不能为空" dc:"清洗任务ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *DataCleanTaskStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		return gerror.New("清洗任务ID不能为空")
	}
	if !validate.InSlice(consts.StatusSlice, in.Status) || in.Status <= 0 {
		return gerror.New("状态不正确")
	}
	return nil
}

type DataCleanTaskStatusModel struct{}

type DataCleanTaskTestInp struct {
	SourceId     int64       `json:"sourceId"     dc:"输入连接ID"`
	SourceConfig *gjson.Json `json:"sourceConfig" dc:"输入配置"`
	CleanConfig  *gjson.Json `json:"cleanConfig"  dc:"清洗配置"`
	SinkConfig   *gjson.Json `json:"sinkConfig"   dc:"输出配置"`
}

func (in *DataCleanTaskTestInp) Filter(ctx context.Context) (err error) {
	if in.SourceId <= 0 {
		return gerror.New("请选择输入数据源")
	}
	return nil
}

type DataCleanTaskTestModel struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"测试结果"`
}

type DataCleanTaskSampleInp struct {
	Id             int64       `json:"id"           dc:"清洗任务ID"`
	SourceId       int64       `json:"sourceId"     dc:"输入连接ID"`
	SourceConfig   *gjson.Json `json:"sourceConfig" dc:"输入配置"`
	CleanConfig    *gjson.Json `json:"cleanConfig"  dc:"清洗配置"`
	Payload        string      `json:"payload"      dc:"样例数据"`
	Limit          int         `json:"limit"        dc:"采样数量"`
	TimeoutSeconds int         `json:"timeoutSeconds" dc:"采样超时秒数"`
}

func (in *DataCleanTaskSampleInp) Filter(ctx context.Context) (err error) {
	if in.SourceId <= 0 {
		return gerror.New("请选择输入数据源")
	}
	if in.Limit <= 0 {
		in.Limit = 50
	}
	if in.Limit > 500 {
		in.Limit = 500
	}
	if in.TimeoutSeconds <= 0 {
		in.TimeoutSeconds = 10
	}
	if in.TimeoutSeconds > 60 {
		in.TimeoutSeconds = 60
	}
	return nil
}

type DataCleanTaskSampleModel struct {
	List  []*DataFieldListModel `json:"list"  dc:"字段列表"`
	Total int                   `json:"total" dc:"字段数量"`
}

type DataFieldListInp struct {
	TaskId int64 `json:"taskId" v:"required#清洗任务ID不能为空" dc:"清洗任务ID"`
}

func (in *DataFieldListInp) Filter(ctx context.Context) (err error) { return }

type DataFieldListModel struct {
	Id           int64       `json:"id"           dc:"字段ID"`
	TaskId       int64       `json:"taskId"       dc:"清洗任务ID"`
	ConnectorId  int64       `json:"connectorId"  dc:"输入连接ID"`
	FieldPath    string      `json:"fieldPath"    dc:"字段路径"`
	FieldName    string      `json:"fieldName"    dc:"字段名称"`
	FieldType    string      `json:"fieldType"    dc:"字段类型"`
	FieldTypes   *gjson.Json `json:"fieldTypes"   dc:"出现过的字段类型集合"`
	SampleValue  *gjson.Json `json:"sampleValue"  dc:"样例值"`
	SampleValues []string    `json:"sampleValues" dc:"样例值列表"`
	Status       int         `json:"status"       dc:"状态"`
	Remark       string      `json:"remark"       dc:"备注"`
	FirstSeenAt  *gtime.Time `json:"firstSeenAt"  dc:"首次发现时间"`
	LastSeenAt   *gtime.Time `json:"lastSeenAt"   dc:"最后发现时间"`
}
