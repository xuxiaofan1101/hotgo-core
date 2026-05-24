package dataclean

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta `path:"/dataClean/list" method:"get" tags:"数据清洗" summary:"获取数据清洗任务列表"`
	sysin.DataCleanTaskListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.DataCleanTaskListModel `json:"list" dc:"数据列表"`
}

type ViewReq struct {
	g.Meta `path:"/dataClean/view" method:"get" tags:"数据清洗" summary:"获取数据清洗任务详情"`
	sysin.DataCleanTaskViewInp
}

type ViewRes struct {
	*sysin.DataCleanTaskViewModel
}

type EditReq struct {
	g.Meta `path:"/dataClean/edit" method:"post" tags:"数据清洗" summary:"新增/编辑数据清洗任务"`
	sysin.DataCleanTaskEditInp
}

type EditRes struct{}

type DeleteReq struct {
	g.Meta `path:"/dataClean/delete" method:"post" tags:"数据清洗" summary:"删除数据清洗任务"`
	sysin.DataCleanTaskDeleteInp
}

type DeleteRes struct{}

type StatusReq struct {
	g.Meta `path:"/dataClean/status" method:"post" tags:"数据清洗" summary:"更新数据清洗任务状态"`
	sysin.DataCleanTaskStatusInp
}

type StatusRes struct{}

type TestReq struct {
	g.Meta `path:"/dataClean/test" method:"post" tags:"数据清洗" summary:"测试数据清洗配置"`
	sysin.DataCleanTaskTestInp
}

type TestRes struct {
	*sysin.DataCleanTaskTestModel
}

type SampleReq struct {
	g.Meta `path:"/dataClean/sample" method:"post" tags:"数据清洗" summary:"采样数据清洗字段"`
	sysin.DataCleanTaskSampleInp
}

type SampleRes struct {
	*sysin.DataCleanTaskSampleModel
}

type FieldListReq struct {
	g.Meta `path:"/dataClean/fieldList" method:"get" tags:"数据清洗" summary:"获取清洗任务采集字段列表"`
	sysin.DataFieldListInp
}

type FieldListRes struct {
	List  []*sysin.DataFieldListModel `json:"list"  dc:"字段列表"`
	Total int                         `json:"total" dc:"字段数量"`
}
