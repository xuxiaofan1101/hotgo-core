package datafieldtemplate

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta `path:"/dataFieldTemplate/list" method:"get" tags:"字段模板" summary:"获取字段模板列表"`
	sysin.DataFieldTemplateListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.DataFieldTemplateListModel `json:"list" dc:"数据列表"`
}

type ViewReq struct {
	g.Meta `path:"/dataFieldTemplate/view" method:"get" tags:"字段模板" summary:"获取字段模板详情"`
	sysin.DataFieldTemplateViewInp
}

type ViewRes struct {
	*sysin.DataFieldTemplateViewModel
}

type EditReq struct {
	g.Meta `path:"/dataFieldTemplate/edit" method:"post" tags:"字段模板" summary:"新增/编辑字段模板"`
	sysin.DataFieldTemplateEditInp
}

type EditRes struct{}

type DeleteReq struct {
	g.Meta `path:"/dataFieldTemplate/delete" method:"post" tags:"字段模板" summary:"删除字段模板"`
	sysin.DataFieldTemplateDeleteInp
}

type DeleteRes struct{}

type StatusReq struct {
	g.Meta `path:"/dataFieldTemplate/status" method:"post" tags:"字段模板" summary:"更新字段模板状态"`
	sysin.DataFieldTemplateStatusInp
}

type StatusRes struct{}
