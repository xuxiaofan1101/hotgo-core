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

type DataFieldTemplateUpdateFields struct {
	Name      string      `json:"name"      dc:"模板名称"`
	Code      string      `json:"code"      dc:"模板编码"`
	EventType string      `json:"eventType" dc:"事件类型"`
	Fields    *gjson.Json `json:"fields"    dc:"标准字段定义"`
	Examples  *gjson.Json `json:"examples"  dc:"标准样例数据"`
	Status    int         `json:"status"    dc:"状态"`
	Remark    string      `json:"remark"    dc:"备注"`
	UpdatedBy int64       `json:"updatedBy" dc:"更新者"`
}

type DataFieldTemplateInsertFields struct {
	Name      string      `json:"name"      dc:"模板名称"`
	Code      string      `json:"code"      dc:"模板编码"`
	EventType string      `json:"eventType" dc:"事件类型"`
	Fields    *gjson.Json `json:"fields"    dc:"标准字段定义"`
	Examples  *gjson.Json `json:"examples"  dc:"标准样例数据"`
	Status    int         `json:"status"    dc:"状态"`
	Remark    string      `json:"remark"    dc:"备注"`
	CreatedBy int64       `json:"createdBy" dc:"创建者"`
}

type DataFieldTemplateEditInp struct {
	entity.DataFieldTemplate
}

func (in *DataFieldTemplateEditInp) Filter(ctx context.Context) (err error) {
	if e := g.Validator().Rules("required").Data(in.Name).Messages("模板名称不能为空").Run(ctx); e != nil {
		return e.Current()
	}
	if e := g.Validator().Rules("required").Data(in.Code).Messages("模板编码不能为空").Run(ctx); e != nil {
		return e.Current()
	}
	if e := g.Validator().Rules("required").Data(in.EventType).Messages("事件类型不能为空").Run(ctx); e != nil {
		return e.Current()
	}
	if in.Status <= 0 {
		in.Status = consts.StatusEnabled
	}
	if !validate.InSlice(consts.StatusSlice, in.Status) {
		return gerror.New("状态不正确")
	}
	return nil
}

type DataFieldTemplateEditModel struct{}

type DataFieldTemplateDeleteInp struct {
	Id interface{} `json:"id" v:"required#字段模板ID不能为空" dc:"字段模板ID"`
}

func (in *DataFieldTemplateDeleteInp) Filter(ctx context.Context) (err error) { return }

type DataFieldTemplateDeleteModel struct{}

type DataFieldTemplateViewInp struct {
	Id int64 `json:"id" v:"required#字段模板ID不能为空" dc:"字段模板ID"`
}

func (in *DataFieldTemplateViewInp) Filter(ctx context.Context) (err error) { return }

type DataFieldTemplateViewModel struct {
	entity.DataFieldTemplate
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma" dc:"创建者摘要信息"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma" dc:"更新者摘要信息"`
}

type DataFieldTemplateListInp struct {
	form.PageReq
	Id        int64         `json:"id"        dc:"字段模板ID"`
	Keyword   string        `json:"keyword"   dc:"关键词"`
	Name      string        `json:"name"      dc:"模板名称"`
	Code      string        `json:"code"      dc:"模板编码"`
	EventType string        `json:"eventType" dc:"事件类型"`
	Status    int           `json:"status"    dc:"状态"`
	CreatedBy string        `json:"createdBy" dc:"创建者"`
	CreatedAt []*gtime.Time `json:"createdAt" dc:"创建时间"`
}

func (in *DataFieldTemplateListInp) Filter(ctx context.Context) (err error) { return }

type DataFieldTemplateListModel struct {
	Id             int64             `json:"id"             dc:"字段模板ID"`
	Name           string            `json:"name"           dc:"模板名称"`
	Code           string            `json:"code"           dc:"模板编码"`
	EventType      string            `json:"eventType"      dc:"事件类型"`
	Fields         *gjson.Json       `json:"fields"         dc:"标准字段定义"`
	Examples       *gjson.Json       `json:"examples"       dc:"标准样例数据"`
	FieldCount     int               `json:"fieldCount"     dc:"字段数量"`
	Status         int               `json:"status"         dc:"状态"`
	Remark         string            `json:"remark"         dc:"备注"`
	CreatedBy      int64             `json:"createdBy"      dc:"创建者"`
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma" dc:"创建者摘要信息"`
	UpdatedBy      int64             `json:"updatedBy"      dc:"更新者"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma" dc:"更新者摘要信息"`
	CreatedAt      *gtime.Time       `json:"createdAt"      dc:"创建时间"`
	UpdatedAt      *gtime.Time       `json:"updatedAt"      dc:"修改时间"`
}

type DataFieldTemplateStatusInp struct {
	Id     int64 `json:"id" v:"required#字段模板ID不能为空" dc:"字段模板ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *DataFieldTemplateStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		return gerror.New("字段模板ID不能为空")
	}
	if !validate.InSlice(consts.StatusSlice, in.Status) || in.Status <= 0 {
		return gerror.New("状态不正确")
	}
	return nil
}

type DataFieldTemplateStatusModel struct{}
