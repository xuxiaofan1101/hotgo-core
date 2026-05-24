package sys

import (
	"context"
	"strings"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/model/do"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSysDataFieldTemplate struct{}

func NewSysDataFieldTemplate() *sSysDataFieldTemplate {
	return &sSysDataFieldTemplate{}
}

func init() {
	service.RegisterSysDataFieldTemplate(NewSysDataFieldTemplate())
}

func (s *sSysDataFieldTemplate) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.DataFieldTemplate.Ctx(ctx), option...)
}

func (s *sSysDataFieldTemplate) List(ctx context.Context, in *sysin.DataFieldTemplateListInp) (list []*sysin.DataFieldTemplateListModel, totalCount int, err error) {
	mod := s.Model(ctx).FieldsPrefix(dao.DataFieldTemplate.Table(), sysin.DataFieldTemplateListModel{})
	cols := dao.DataFieldTemplate.Columns()

	if in.Id > 0 {
		mod = mod.Where(cols.Id, in.Id)
	}
	if in.Keyword != "" {
		keyword := "%" + strings.TrimSpace(in.Keyword) + "%"
		mod = mod.Where("(`"+cols.Name+"` LIKE ? OR `"+cols.Code+"` LIKE ? OR `"+cols.EventType+"` LIKE ?)", keyword, keyword, keyword)
	}
	if in.Name != "" {
		mod = mod.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.Code != "" {
		mod = mod.WhereLike(cols.Code, "%"+in.Code+"%")
	}
	if in.EventType != "" {
		mod = mod.Where(cols.EventType, in.EventType)
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}
	if in.CreatedBy != "" {
		ids, err := service.AdminMember().GetIdsByKeyword(ctx, in.CreatedBy)
		if err != nil {
			return nil, 0, err
		}
		mod = mod.WhereIn(cols.CreatedBy, ids)
	}
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(cols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	mod = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Hook(hook.MemberSummary)
	if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
		return nil, 0, gerror.Wrap(err, "获取字段模板列表失败，请稍后重试！")
	}
	for _, item := range list {
		if item.Fields != nil {
			item.FieldCount = len(item.Fields.Array())
		}
	}
	return
}

func (s *sSysDataFieldTemplate) Edit(ctx context.Context, in *sysin.DataFieldTemplateEditInp) (err error) {
	if err = validateDataFieldTemplateFields(in.Fields); err != nil {
		return err
	}
	if err = hgorm.IsUnique(ctx, &dao.DataFieldTemplate, g.Map{dao.DataFieldTemplate.Columns().Code: in.Code}, "模板编码已存在", in.Id); err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		if in.Id > 0 {
			in.UpdatedBy = contexts.GetUserId(ctx)
			if _, err = s.Model(ctx).
				Fields(sysin.DataFieldTemplateUpdateFields{}).
				WherePri(in.Id).Data(in).Update(); err != nil {
				return gerror.Wrap(err, "修改字段模板失败，请稍后重试！")
			}
			return nil
		}
		in.CreatedBy = contexts.GetUserId(ctx)
		if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
			Fields(sysin.DataFieldTemplateInsertFields{}).
			Data(in).OmitEmptyData().Insert(); err != nil {
			return gerror.Wrap(err, "新增字段模板失败，请稍后重试！")
		}
		return nil
	})
}

func (s *sSysDataFieldTemplate) Delete(ctx context.Context, in *sysin.DataFieldTemplateDeleteInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataFieldTemplate{
		DeletedBy: contexts.GetUserId(ctx),
		DeletedAt: gtime.Now(),
	}).Update(); err != nil {
		return gerror.Wrap(err, "删除字段模板失败，请稍后重试！")
	}
	return nil
}

func (s *sSysDataFieldTemplate) View(ctx context.Context, in *sysin.DataFieldTemplateViewInp) (res *sysin.DataFieldTemplateViewModel, err error) {
	if err = s.Model(ctx).WherePri(in.Id).Hook(hook.MemberSummary).Scan(&res); err != nil {
		return nil, gerror.Wrap(err, "获取字段模板信息失败，请稍后重试！")
	}
	return
}

func (s *sSysDataFieldTemplate) Status(ctx context.Context, in *sysin.DataFieldTemplateStatusInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataFieldTemplate{
		Status:    in.Status,
		UpdatedBy: contexts.GetUserId(ctx),
	}).Update(); err != nil {
		return gerror.Wrap(err, "更新字段模板状态失败，请稍后重试！")
	}
	return nil
}

func validateDataFieldTemplateFields(fields *gjson.Json) error {
	if fields == nil || len(fields.Array()) == 0 {
		return gerror.New("请配置标准字段")
	}
	return nil
}
